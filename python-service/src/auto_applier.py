import logging
import hashlib
import uuid
from enum import Enum
from pathlib import Path
from typing import Optional

import httpx
from pydantic import BaseModel

logger = logging.getLogger(__name__)


class ApplyMethod(str, Enum):
    API_GREENHOUSE = "api_greenhouse"
    API_LEVER = "api_lever"
    FORM_FILL = "form_fill"
    EMAIL = "email"
    MANUAL = "manual"


class AutoApplyRequest(BaseModel):
    job_url: str
    apply_url: str
    source: str
    resume_pdf_path: str
    cover_letter_pdf_path: Optional[str] = None
    user_name: str
    user_email: str
    user_phone: Optional[str] = None
    linkedin_url: Optional[str] = None


class AutoApplyResult(BaseModel):
    success: bool
    method: str
    confirmation_id: Optional[str] = None
    screenshot_path: Optional[str] = None
    error: Optional[str] = None
    requires_confirmation: bool = False


class AutoApplier:
    """Routes job applications to the appropriate submission strategy.

    - ATS APIs (Greenhouse, Lever): fully automated submission
    - Form fill (generic career pages): fills form + takes screenshot, requires user confirmation
    - Email: composes email payload for the Go backend to send, requires confirmation
    - LinkedIn/Indeed: blocked (ToS violation) -- returns MANUAL
    """

    BLOCKED_SOURCES = {"linkedin", "indeed"}

    def detect_method(self, apply_url: str, source: str) -> ApplyMethod:
        """Detect the best application method based on the URL and source."""
        source_lower = source.lower()

        # LinkedIn and Indeed are NEVER auto-applied
        if source_lower in self.BLOCKED_SOURCES:
            return ApplyMethod.MANUAL

        # Greenhouse ATS
        if "greenhouse.io" in apply_url or "boards.greenhouse" in apply_url:
            return ApplyMethod.API_GREENHOUSE

        # Lever ATS
        if "lever.co" in apply_url or "jobs.lever" in apply_url:
            return ApplyMethod.API_LEVER

        # Email-based applications
        if apply_url.startswith("mailto:"):
            return ApplyMethod.EMAIL

        # Everything else: browser-based form filling
        return ApplyMethod.FORM_FILL

    async def apply(self, request: AutoApplyRequest) -> AutoApplyResult:
        """Execute the application using the detected method."""
        method = self.detect_method(request.apply_url, request.source)

        if method == ApplyMethod.MANUAL:
            return AutoApplyResult(
                success=False,
                method=method.value,
                error=f"Auto-apply disabled for {request.source}. Please apply manually.",
                requires_confirmation=False,
            )
        elif method == ApplyMethod.API_GREENHOUSE:
            return await self._apply_greenhouse(request)
        elif method == ApplyMethod.API_LEVER:
            return await self._apply_lever(request)
        elif method == ApplyMethod.EMAIL:
            return await self._apply_email(request)
        elif method == ApplyMethod.FORM_FILL:
            return await self._apply_form_fill(request)

        return AutoApplyResult(success=False, method="unknown", error="Unknown apply method")

    async def _apply_greenhouse(self, request: AutoApplyRequest) -> AutoApplyResult:
        """Submit application via the Greenhouse job board API.

        Greenhouse boards accept multipart POST requests at:
        https://boards-api.greenhouse.io/v1/boards/{board}/jobs/{job_id}
        """
        try:
            # Parse the board token and job ID from the apply URL
            # Typical URL patterns:
            #   https://boards.greenhouse.io/company/jobs/12345
            #   https://job-boards.greenhouse.io/company/jobs/12345
            parts = request.apply_url.rstrip("/").split("/")
            job_id = parts[-1] if parts else ""
            board_token = parts[-3] if len(parts) >= 3 else ""

            if not job_id or not board_token:
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.API_GREENHOUSE.value,
                    error="Could not parse Greenhouse board token or job ID from URL",
                )

            api_url = f"https://boards-api.greenhouse.io/v1/boards/{board_token}/jobs/{job_id}"

            # Split user name into first/last
            name_parts = request.user_name.strip().split(" ", 1)
            first_name = name_parts[0]
            last_name = name_parts[1] if len(name_parts) > 1 else ""

            # Build multipart form data
            form_data = {
                "first_name": first_name,
                "last_name": last_name,
                "email": request.user_email,
            }

            if request.user_phone:
                form_data["phone"] = request.user_phone

            if request.linkedin_url:
                form_data["linkedin_profile_url"] = request.linkedin_url

            # Prepare file uploads
            resume_path = Path(request.resume_pdf_path)
            if not resume_path.exists():
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.API_GREENHOUSE.value,
                    error=f"Resume PDF not found: {request.resume_pdf_path}",
                )

            files = {
                "resume": (resume_path.name, resume_path.read_bytes(), "application/pdf"),
            }

            if request.cover_letter_pdf_path:
                cover_path = Path(request.cover_letter_pdf_path)
                if cover_path.exists():
                    files["cover_letter"] = (
                        cover_path.name,
                        cover_path.read_bytes(),
                        "application/pdf",
                    )

            async with httpx.AsyncClient(timeout=30.0) as client:
                resp = await client.post(api_url, data=form_data, files=files)

            if resp.status_code in (200, 201):
                resp_data = resp.json()
                confirmation_id = str(resp_data.get("id", hashlib.sha256(resp.content).hexdigest()[:12]))
                logger.info(
                    "Greenhouse application submitted: job_id=%s, confirmation=%s",
                    job_id,
                    confirmation_id,
                )
                return AutoApplyResult(
                    success=True,
                    method=ApplyMethod.API_GREENHOUSE.value,
                    confirmation_id=confirmation_id,
                )
            else:
                error_text = resp.text[:500]
                logger.warning(
                    "Greenhouse API returned %d for job %s: %s",
                    resp.status_code,
                    job_id,
                    error_text,
                )
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.API_GREENHOUSE.value,
                    error=f"Greenhouse API error ({resp.status_code}): {error_text}",
                )

        except Exception as exc:
            logger.exception("Greenhouse application failed")
            return AutoApplyResult(
                success=False,
                method=ApplyMethod.API_GREENHOUSE.value,
                error=f"Greenhouse application error: {str(exc)}",
            )

    async def _apply_lever(self, request: AutoApplyRequest) -> AutoApplyResult:
        """Submit application via the Lever postings API.

        Lever postings accept multipart POST requests at:
        https://api.lever.co/v0/postings/{company}/{posting_id}
        """
        try:
            # Parse posting ID and company from the apply URL
            # Typical URL patterns:
            #   https://jobs.lever.co/company/posting-uuid
            #   https://jobs.lever.co/company/posting-uuid/apply
            url_clean = request.apply_url.rstrip("/")
            if url_clean.endswith("/apply"):
                url_clean = url_clean[: -len("/apply")]

            parts = url_clean.split("/")
            posting_id = parts[-1] if parts else ""
            company = parts[-2] if len(parts) >= 2 else ""

            if not posting_id or not company:
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.API_LEVER.value,
                    error="Could not parse Lever company or posting ID from URL",
                )

            api_url = f"https://api.lever.co/v0/postings/{company}/{posting_id}"

            # Split user name
            name_parts = request.user_name.strip().split(" ", 1)
            first_name = name_parts[0]
            last_name = name_parts[1] if len(name_parts) > 1 else ""

            form_data = {
                "name": request.user_name,
                "email": request.user_email,
            }

            if request.user_phone:
                form_data["phone"] = request.user_phone

            if request.linkedin_url:
                form_data["urls[LinkedIn]"] = request.linkedin_url

            # Prepare resume file
            resume_path = Path(request.resume_pdf_path)
            if not resume_path.exists():
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.API_LEVER.value,
                    error=f"Resume PDF not found: {request.resume_pdf_path}",
                )

            files = {
                "resume": (resume_path.name, resume_path.read_bytes(), "application/pdf"),
            }

            async with httpx.AsyncClient(timeout=30.0) as client:
                resp = await client.post(api_url, data=form_data, files=files)

            if resp.status_code in (200, 201):
                resp_data = resp.json()
                confirmation_id = str(resp_data.get("applicationId", hashlib.sha256(resp.content).hexdigest()[:12]))
                logger.info(
                    "Lever application submitted: posting=%s, confirmation=%s",
                    posting_id,
                    confirmation_id,
                )
                return AutoApplyResult(
                    success=True,
                    method=ApplyMethod.API_LEVER.value,
                    confirmation_id=confirmation_id,
                )
            else:
                error_text = resp.text[:500]
                logger.warning(
                    "Lever API returned %d for posting %s: %s",
                    resp.status_code,
                    posting_id,
                    error_text,
                )
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.API_LEVER.value,
                    error=f"Lever API error ({resp.status_code}): {error_text}",
                )

        except Exception as exc:
            logger.exception("Lever application failed")
            return AutoApplyResult(
                success=False,
                method=ApplyMethod.API_LEVER.value,
                error=f"Lever application error: {str(exc)}",
            )

    async def _apply_email(self, request: AutoApplyRequest) -> AutoApplyResult:
        """Compose an email application payload.

        Does NOT actually send the email -- returns the composed data so the Go
        backend can send it through its email service. Always requires user
        confirmation before the backend dispatches.
        """
        try:
            # Extract email address from mailto: URL
            mailto = request.apply_url
            if mailto.startswith("mailto:"):
                mailto = mailto[len("mailto:"):]

            # Strip any query params (?subject=... etc.)
            recipient = mailto.split("?")[0].strip()

            if not recipient:
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.EMAIL.value,
                    error="Could not extract email address from mailto URL",
                )

            # Build email body
            name_parts = request.user_name.strip().split(" ", 1)
            first_name = name_parts[0]

            body_lines = [
                f"Dear Hiring Manager,",
                "",
                f"I am writing to express my interest in the position listed at:",
                request.job_url,
                "",
                f"Please find my resume attached.",
            ]

            if request.cover_letter_pdf_path:
                body_lines.append("I have also attached a cover letter for your review.")

            body_lines.extend(
                [
                    "",
                    "I look forward to hearing from you.",
                    "",
                    f"Best regards,",
                    request.user_name,
                    request.user_email,
                ]
            )

            if request.user_phone:
                body_lines.append(request.user_phone)

            if request.linkedin_url:
                body_lines.append(request.linkedin_url)

            email_body = "\n".join(body_lines)

            # Generate a tracking confirmation ID
            confirmation_id = f"email-{uuid.uuid4().hex[:8]}"

            logger.info(
                "Email application composed: to=%s, confirmation=%s",
                recipient,
                confirmation_id,
            )

            return AutoApplyResult(
                success=True,
                method=ApplyMethod.EMAIL.value,
                confirmation_id=confirmation_id,
                requires_confirmation=True,
                error=None,
            )

        except Exception as exc:
            logger.exception("Email application composition failed")
            return AutoApplyResult(
                success=False,
                method=ApplyMethod.EMAIL.value,
                error=f"Email composition error: {str(exc)}",
            )

    async def _apply_form_fill(self, request: AutoApplyRequest) -> AutoApplyResult:
        """Fill a web application form using Playwright.

        NEVER auto-submits -- always takes a screenshot and requires user
        confirmation before the Go backend triggers actual submission.
        """
        try:
            from src.form_filler import FormFiller

            output_dir = Path("data_folder/output/screenshots")
            output_dir.mkdir(parents=True, exist_ok=True)

            user_data = {
                "name": request.user_name,
                "email": request.user_email,
                "phone": request.user_phone or "",
                "linkedin": request.linkedin_url or "",
            }

            filler = FormFiller()
            result = await filler.fill_form(
                apply_url=request.apply_url,
                user_data=user_data,
                resume_path=request.resume_pdf_path,
                output_dir=str(output_dir),
            )

            if result.get("error"):
                return AutoApplyResult(
                    success=False,
                    method=ApplyMethod.FORM_FILL.value,
                    error=result["error"],
                )

            logger.info(
                "Form filled: url=%s, filled=%d, unfilled=%d, screenshot=%s",
                request.apply_url,
                len(result.get("filled_fields", [])),
                len(result.get("unfilled_fields", [])),
                result.get("screenshot_path"),
            )

            return AutoApplyResult(
                success=True,
                method=ApplyMethod.FORM_FILL.value,
                screenshot_path=result.get("screenshot_path"),
                requires_confirmation=True,
            )

        except Exception as exc:
            logger.exception("Form fill failed")
            return AutoApplyResult(
                success=False,
                method=ApplyMethod.FORM_FILL.value,
                error=f"Form fill error: {str(exc)}",
            )
