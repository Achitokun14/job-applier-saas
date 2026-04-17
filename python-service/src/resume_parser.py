import logging
from typing import Dict, List, Optional
from pydantic import BaseModel, Field

logger = logging.getLogger(__name__)


class ParsedResume(BaseModel):
    name: str = ""
    email: str = ""
    phone: str = ""
    location: str = ""
    summary: str = ""
    experience: List[dict] = Field(default_factory=list)  # [{title, company, dates, description}]
    education: List[dict] = Field(default_factory=list)  # [{degree, school, dates}]
    skills: List[str] = Field(default_factory=list)
    projects: List[dict] = Field(default_factory=list)  # [{name, description}]
    certifications: List[str] = Field(default_factory=list)
    languages: List[str] = Field(default_factory=list)


class ResumeParser:
    def parse_pdf(
        self,
        pdf_path: str,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Extract structured resume data from a PDF file."""
        # Step 1: Extract text using PyMuPDF4LLM (outputs clean markdown)
        try:
            import pymupdf4llm

            markdown_text = pymupdf4llm.to_markdown(pdf_path)
        except Exception as e:
            logger.warning(f"PyMuPDF4LLM failed: {e}, falling back to pdfplumber")
            import pdfplumber

            with pdfplumber.open(pdf_path) as pdf:
                markdown_text = "\n".join(
                    page.extract_text() or "" for page in pdf.pages
                )

        if not markdown_text.strip():
            return ParsedResume().model_dump()

        # Step 2: Use LLM + Instructor to parse into structured data
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )
            parsed = client.chat.completions.create(
                model=model,
                response_model=ParsedResume,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "Extract structured resume data from the following text. "
                            "Be thorough and accurate. For experience entries, extract "
                            "title, company, dates, and description. For education, "
                            "extract degree, school, and dates."
                        ),
                    },
                    {"role": "user", "content": markdown_text[:8000]},
                ],
            )
            return parsed.model_dump()
        except Exception as e:
            logger.error(f"LLM parsing failed: {e}")
            return {"raw_text": markdown_text, "parse_error": str(e)}
