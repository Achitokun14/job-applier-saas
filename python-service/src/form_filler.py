import logging
import uuid
from pathlib import Path
from typing import Dict, List, Optional

logger = logging.getLogger(__name__)


# Maps logical field names to lists of CSS selectors commonly used by career pages.
FIELD_SELECTORS: Dict[str, List[str]] = {
    "name": [
        'input[name*="name" i]',
        'input[id*="name" i]',
        'input[placeholder*="name" i]',
        'input[aria-label*="name" i]',
        'input[autocomplete="name"]',
    ],
    "first_name": [
        'input[name*="first" i]',
        'input[id*="first" i]',
        'input[placeholder*="first" i]',
        'input[aria-label*="first name" i]',
        'input[autocomplete="given-name"]',
    ],
    "last_name": [
        'input[name*="last" i]',
        'input[id*="last" i]',
        'input[placeholder*="last" i]',
        'input[aria-label*="last name" i]',
        'input[autocomplete="family-name"]',
    ],
    "email": [
        'input[type="email"]',
        'input[name*="email" i]',
        'input[id*="email" i]',
        'input[placeholder*="email" i]',
        'input[autocomplete="email"]',
    ],
    "phone": [
        'input[type="tel"]',
        'input[name*="phone" i]',
        'input[id*="phone" i]',
        'input[placeholder*="phone" i]',
        'input[autocomplete="tel"]',
    ],
    "linkedin": [
        'input[name*="linkedin" i]',
        'input[id*="linkedin" i]',
        'input[placeholder*="linkedin" i]',
        'input[aria-label*="linkedin" i]',
    ],
}

# Selectors for file upload inputs (resume)
RESUME_UPLOAD_SELECTORS = [
    'input[type="file"][name*="resume" i]',
    'input[type="file"][name*="cv" i]',
    'input[type="file"][id*="resume" i]',
    'input[type="file"][id*="cv" i]',
    'input[type="file"][accept*="pdf"]',
    'input[type="file"]',
]


class FormFiller:
    """Fills web-based job application forms using Playwright.

    Key safety constraints:
    - NEVER clicks submit buttons
    - Always takes a screenshot of the filled form for user review
    - Gracefully skips fields it cannot locate
    """

    async def fill_form(
        self,
        apply_url: str,
        user_data: Dict[str, str],
        resume_path: str,
        output_dir: str,
    ) -> Dict:
        """Navigate to an application form, fill fields, upload resume, and screenshot.

        Args:
            apply_url: URL of the job application form.
            user_data: Dict with keys: name, email, phone, linkedin.
            resume_path: Path to the resume PDF file.
            output_dir: Directory to save screenshots.

        Returns:
            Dict with keys: filled_fields, unfilled_fields, screenshot_path, error.
        """
        filled_fields: List[str] = []
        unfilled_fields: List[str] = []
        screenshot_path: Optional[str] = None

        try:
            from playwright.async_api import async_playwright
        except ImportError:
            return {
                "filled_fields": [],
                "unfilled_fields": [],
                "screenshot_path": None,
                "error": "Playwright is not installed. Run: pip install playwright && playwright install chromium",
            }

        try:
            async with async_playwright() as pw:
                browser = await pw.chromium.launch(headless=True)
                context = await browser.new_context(
                    viewport={"width": 1280, "height": 900},
                    user_agent=(
                        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
                        "AppleWebKit/537.36 (KHTML, like Gecko) "
                        "Chrome/120.0.0.0 Safari/537.36"
                    ),
                )
                page = await context.new_page()

                # Navigate to the application page
                await page.goto(apply_url, wait_until="domcontentloaded", timeout=30000)
                await page.wait_for_timeout(2000)  # Let JS render

                # Take a "before" screenshot
                session_id = uuid.uuid4().hex[:8]
                output_path = Path(output_dir)
                output_path.mkdir(parents=True, exist_ok=True)

                before_path = str(output_path / f"form_before_{session_id}.png")
                await page.screenshot(path=before_path, full_page=True)

                # Derive first/last name from full name
                name_parts = user_data.get("name", "").strip().split(" ", 1)
                first_name = name_parts[0] if name_parts else ""
                last_name = name_parts[1] if len(name_parts) > 1 else ""

                # Map logical fields to their values
                field_values = {
                    "name": user_data.get("name", ""),
                    "first_name": first_name,
                    "last_name": last_name,
                    "email": user_data.get("email", ""),
                    "phone": user_data.get("phone", ""),
                    "linkedin": user_data.get("linkedin", ""),
                }

                # Try to fill each field
                for field_name, selectors in FIELD_SELECTORS.items():
                    value = field_values.get(field_name, "")
                    if not value:
                        continue

                    field_filled = False
                    for selector in selectors:
                        try:
                            element = page.locator(selector).first
                            if await element.count() > 0:
                                is_visible = await element.is_visible()
                                if is_visible:
                                    await element.click()
                                    await element.fill(value)
                                    filled_fields.append(field_name)
                                    field_filled = True
                                    logger.debug(
                                        "Filled field '%s' using selector '%s'",
                                        field_name,
                                        selector,
                                    )
                                    break
                        except Exception:
                            continue

                    if not field_filled and value:
                        unfilled_fields.append(field_name)

                # Try to upload resume
                resume_file = Path(resume_path)
                if resume_file.exists():
                    resume_uploaded = False
                    for selector in RESUME_UPLOAD_SELECTORS:
                        try:
                            file_input = page.locator(selector).first
                            if await file_input.count() > 0:
                                await file_input.set_input_files(str(resume_file))
                                filled_fields.append("resume_upload")
                                resume_uploaded = True
                                logger.debug(
                                    "Uploaded resume using selector '%s'", selector
                                )
                                break
                        except Exception:
                            continue

                    if not resume_uploaded:
                        unfilled_fields.append("resume_upload")
                else:
                    unfilled_fields.append("resume_upload (file not found)")

                # Wait a moment for any client-side validation to render
                await page.wait_for_timeout(1000)

                # Take "after" screenshot -- this is the one sent to the user for review
                after_path = str(output_path / f"form_filled_{session_id}.png")
                await page.screenshot(path=after_path, full_page=True)
                screenshot_path = after_path

                logger.info(
                    "Form fill complete: url=%s, filled=%s, unfilled=%s",
                    apply_url,
                    filled_fields,
                    unfilled_fields,
                )

                # IMPORTANT: Do NOT click any submit button.
                # The browser is closed here. Submission is handled separately
                # after user confirmation.
                await browser.close()

        except Exception as exc:
            logger.exception("Form filler encountered an error")
            return {
                "filled_fields": filled_fields,
                "unfilled_fields": unfilled_fields,
                "screenshot_path": screenshot_path,
                "error": f"Form fill error: {str(exc)}",
            }

        return {
            "filled_fields": filled_fields,
            "unfilled_fields": unfilled_fields,
            "screenshot_path": screenshot_path,
            "error": None,
        }
