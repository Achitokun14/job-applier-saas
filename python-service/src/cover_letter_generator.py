import os
import logging
from datetime import date
from typing import Optional, Dict, Any, List
from pathlib import Path

from pydantic import BaseModel, Field
from jinja2 import Environment, FileSystemLoader

logger = logging.getLogger(__name__)

# Resolve templates directory relative to this file's location
TEMPLATES_DIR = Path(__file__).resolve().parent.parent / "templates"

jinja_env = Environment(
    loader=FileSystemLoader(str(TEMPLATES_DIR)),
    autoescape=True,
)


class CoverLetterContent(BaseModel):
    greeting: str = Field(description="Opening greeting, e.g. 'Dear Hiring Manager,'")
    opening_paragraph: str = Field(description="Opening paragraph expressing interest")
    body_paragraphs: List[str] = Field(description="Body paragraphs highlighting relevant experience")
    closing_paragraph: str = Field(description="Closing paragraph with call to action")
    sign_off: str = Field(description="Sign-off, e.g. 'Sincerely, [Your Name]'")


class CoverLetterGenerator:
    def __init__(self):
        self.llm_api_key = os.getenv("LLM_API_KEY", "")
        self.llm_model = os.getenv("LLM_MODEL", "gpt-4o-mini")
        self.llm_provider = os.getenv("LLM_PROVIDER", "openai")

    def generate(
        self,
        resume_text: str,
        job_description: str,
        job_url: Optional[str] = None,
        company_name: Optional[str] = None,
        job_title: Optional[str] = None,
        output_path: str = "output/cover_letter.pdf",
        llm_provider: Optional[str] = None,
        llm_api_key: Optional[str] = None,
        llm_model: Optional[str] = None,
    ) -> Dict[str, Any]:
        api_key = llm_api_key or self.llm_api_key
        if api_key:
            html_content = self._generate_with_llm(
                resume_text, job_description, company_name, job_title,
                provider=llm_provider or self.llm_provider,
                api_key=api_key,
                model=llm_model or self.llm_model,
            )
        else:
            html_content = self._generate_template(
                resume_text, job_description, company_name, job_title
            )

        word_count = len(html_content.split())

        self._save_html(html_content, output_path.replace('.pdf', '.html'))
        self._generate_pdf(html_content, output_path)

        return {
            "pdf_path": output_path,
            "html_content": html_content,
            "word_count": word_count
        }

    def _generate_with_llm(self, resume_text, job_description, company_name, job_title, provider, api_key, model):
        try:
            from src.llm_client import get_client
            client, full_model = get_client(provider=provider, api_key=api_key, model=model)
            cover_letter = client.chat.completions.create(
                model=full_model, response_model=CoverLetterContent, temperature=0.7,
                messages=[{"role": "user", "content": (
                    "Write a professional cover letter for the following:\n\n"
                    f"Job Title: {job_title or 'the position'}\n"
                    f"Company: {company_name or 'your company'}\n"
                    f"Job Description: {job_description}\n\n"
                    f"Candidate Resume: {resume_text}\n\n"
                    "Write a compelling cover letter that:\n"
                    "1. Opens with enthusiasm for the specific role and company\n"
                    "2. Highlights relevant experience from the resume\n"
                    "3. Shows understanding of the job requirements\n"
                    "4. Closes with a call to action\n\n"
                    "Keep it professional, concise (3-4 paragraphs), and tailored to the job.")}])
            paragraphs = [cover_letter.greeting, cover_letter.opening_paragraph]
            paragraphs.extend(cover_letter.body_paragraphs)
            paragraphs.append(cover_letter.closing_paragraph)
            return self._render_cover_letter_template(paragraphs=paragraphs, company_name=company_name, job_title=job_title, sign_off=cover_letter.sign_off)
        except Exception as e:
            logger.warning(f"LLM generation failed: {e}")
            return self._generate_template(resume_text, job_description, company_name, job_title)

    def _generate_template(self, resume_text, job_description, company_name, job_title):
        company = company_name or "the company"
        position = job_title or "this position"
        paragraphs = [
            f"Dear Hiring Manager,",
            f"I am writing to express my strong interest in the {position} position at {company}. With my background and experience, I believe I would be a valuable addition to your team.",
            self._extract_highlights(resume_text),
            "I am particularly drawn to this opportunity because it aligns well with my career goals and expertise. The job requirements you've outlined match my skill set, and I am confident I can make meaningful contributions from day one.",
            "I would welcome the opportunity to discuss how my experience and skills can benefit your team. Thank you for considering my application.",
        ]
        return self._render_cover_letter_template(paragraphs=paragraphs, company_name=company_name, job_title=job_title, sign_off="Sincerely,")

    def _render_cover_letter_template(self, paragraphs, company_name=None, job_title=None, sign_off="Sincerely,", sender_name=None, sender_email=None, sender_phone=None, sender_location=None):
        template = jinja_env.get_template("cover_letter.html")
        today = date.today().strftime("%B %d, %Y")
        subject = f"Application for {job_title}" if job_title else ""
        return template.render(sender_name=sender_name or "", sender_email=sender_email or "", sender_phone=sender_phone or "", sender_location=sender_location or "", date=today, recipient_name="Hiring Manager", company_name=company_name or "", company_address="", subject=subject, paragraphs=paragraphs, sign_off=sign_off)

    def _extract_highlights(self, resume_text):
        return "Throughout my career, I have developed strong technical and interpersonal skills that have enabled me to deliver results in complex environments. My experience has taught me the value of collaboration, innovation, and continuous learning."

    def _save_html(self, html_content, path):
        Path(path).parent.mkdir(parents=True, exist_ok=True)
        with open(path, 'w', encoding='utf-8') as f:
            f.write(html_content)

    def _generate_pdf(self, html_content, output_path):
        Path(output_path).parent.mkdir(parents=True, exist_ok=True)
        try:
            self._generate_pdf_weasyprint(html_content, output_path)
            logger.info(f"Cover letter PDF generated with WeasyPrint: {output_path}")
            return
        except Exception as e:
            logger.warning(f"WeasyPrint PDF generation failed: {e}")
        try:
            self._generate_pdf_selenium(html_content, output_path)
            logger.info(f"Cover letter PDF generated with Selenium: {output_path}")
            return
        except Exception as e:
            logger.warning(f"Selenium PDF generation failed: {e}")
        try:
            self._generate_pdf_reportlab(html_content, output_path)
            logger.info(f"Cover letter PDF generated with ReportLab: {output_path}")
        except Exception as e:
            logger.error(f"All PDF generation methods failed: {e}")
            raise RuntimeError(f"PDF generation failed: {e}")

    def _generate_pdf_weasyprint(self, html_content, output_path):
        import weasyprint
        weasyprint.HTML(string=html_content).write_pdf(output_path)
        return output_path

    def _generate_pdf_selenium(self, html_content, output_path):
        from selenium import webdriver
        from selenium.webdriver.chrome.options import Options
        html_path = output_path.replace('.pdf', '.html')
        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disable-dev-shm-usage')
        driver = webdriver.Chrome(options=chrome_options)
        try:
            driver.get(f'file://{os.path.abspath(html_path)}')
            driver.print_page(output_path)
        finally:
            driver.quit()

    def _generate_pdf_reportlab(self, html_content, output_path):
        from reportlab.lib.pagesizes import letter
        from reportlab.platypus import SimpleDocTemplate, Paragraph
        from reportlab.lib.styles import getSampleStyleSheet
        doc = SimpleDocTemplate(output_path, pagesize=letter)
        styles = getSampleStyleSheet()
        story = [Paragraph("Cover Letter PDF - See HTML version for full formatting", styles['Normal'])]
        doc.build(story)
