import os
import yaml
import json
import logging
from typing import Optional, Dict, Any, List
from pathlib import Path

from pydantic import BaseModel, Field
from jinja2 import Environment, FileSystemLoader

logger = logging.getLogger(__name__)

TEMPLATES_DIR = Path(__file__).resolve().parent.parent / "templates"

jinja_env = Environment(
    loader=FileSystemLoader(str(TEMPLATES_DIR)),
    autoescape=True,
)

STYLE_TEMPLATE_MAP = {
    "modern": "resume_modern.html",
    "classic": "resume_classic.html",
    "minimal": "resume_minimal.html",
    "creative": "resume_modern.html",
    "professional": "resume_classic.html",
}


class ResumeTailoringSuggestions(BaseModel):
    skills_to_emphasize: List[str] = Field(default_factory=list, description="Skills from the resume to emphasize for this job")
    experience_rewrites: Dict[str, str] = Field(default_factory=dict, description="Mapping of original experience bullet to rewritten version")
    keywords_to_add: List[str] = Field(default_factory=list, description="Keywords from the job description to weave into the resume")
    summary_suggestion: Optional[str] = Field(default=None, description="Suggested rewrite of the professional summary")


class ResumeGenerator:
    def __init__(self):
        self.llm_api_key = os.getenv("LLM_API_KEY", "")
        self.llm_model = os.getenv("LLM_MODEL", "gpt-4o-mini")
        self.llm_provider = os.getenv("LLM_PROVIDER", "openai")
        self.styles_dir = Path(__file__).parent / "styles"

    def generate(self, resume_yaml, style="modern", job_description=None, output_path="output/resume.pdf", llm_provider=None, llm_api_key=None, llm_model=None):
        try:
            resume_data = yaml.safe_load(resume_yaml)
        except yaml.YAMLError:
            raise ValueError("Invalid YAML format for resume")
        if job_description:
            resume_data = self._tailor_to_job(resume_data, job_description, provider=llm_provider or self.llm_provider, api_key=llm_api_key or self.llm_api_key, model=llm_model or self.llm_model)
        html_content = self._generate_html(resume_data, style)
        word_count = len(html_content.split())
        self._save_html(html_content, output_path.replace('.pdf', '.html'))
        self._generate_pdf(html_content, output_path)
        return {"pdf_path": output_path, "html_content": html_content, "word_count": word_count, "style": style}

    def _tailor_to_job(self, resume_data, job_description, provider, api_key, model):
        if not api_key:
            return resume_data
        try:
            from src.llm_client import get_client
            client, full_model = get_client(provider=provider, api_key=api_key, model=model)
            suggestions = client.chat.completions.create(
                model=full_model, response_model=ResumeTailoringSuggestions, temperature=0.7,
                messages=[{"role": "user", "content": (f"Given this resume data:\n{json.dumps(resume_data, indent=2)}\n\nAnd this job description:\n{job_description}\n\nSuggest improvements to tailor the resume for this job. Include skills to emphasize, experience bullet rewrites, keywords to add, and a summary suggestion.")}])
            if suggestions.summary_suggestion:
                personal = resume_data.get("personal_information", {})
                personal["summary"] = suggestions.summary_suggestion
                resume_data["personal_information"] = personal
            if suggestions.skills_to_emphasize:
                existing_skills = resume_data.get("skills", [])
                for skill in suggestions.skills_to_emphasize:
                    if skill not in existing_skills:
                        existing_skills.insert(0, skill)
                resume_data["skills"] = existing_skills
            if suggestions.experience_rewrites:
                for exp in resume_data.get("experience_details", []):
                    new_responsibilities = []
                    for resp in exp.get("responsibilities", []):
                        new_responsibilities.append(suggestions.experience_rewrites.get(resp, resp))
                    exp["responsibilities"] = new_responsibilities
            return resume_data
        except Exception as e:
            logger.warning(f"LLM tailoring failed: {e}")
            return resume_data

    def _generate_html(self, resume_data, style):
        personal = resume_data.get("personal_information", {})
        education_raw = resume_data.get("education_details", [])
        experience_raw = resume_data.get("experience_details", [])
        skills_raw = resume_data.get("skills", [])
        projects_raw = resume_data.get("projects", [])
        skills = []
        for skill in skills_raw:
            if isinstance(skill, str):
                skills.append(skill)
            elif isinstance(skill, dict):
                skills.append(skill.get("name", ""))
        experience = []
        for exp in experience_raw:
            experience.append({"title": exp.get("position", exp.get("title", "")), "company": exp.get("company", ""), "location": exp.get("location", ""), "dates": f"{exp.get('start_date', '')} - {exp.get('end_date', 'Present')}", "description": exp.get("description", ""), "responsibilities": exp.get("responsibilities", [])})
        education = []
        for edu in education_raw:
            education.append({"degree": edu.get("degree", ""), "institution": edu.get("institution", ""), "dates": edu.get("graduation_date", edu.get("dates", "")), "description": edu.get("description", "")})
        projects = []
        for proj in projects_raw:
            projects.append({"name": proj.get("name", ""), "description": proj.get("description", "")})
        template_file = STYLE_TEMPLATE_MAP.get(style, "resume_modern.html")
        template = jinja_env.get_template(template_file)
        return template.render(name=personal.get("name", "Your Name"), email=personal.get("email", ""), phone=personal.get("phone", ""), location=personal.get("location", ""), linkedin=personal.get("linkedin", ""), github=personal.get("github", ""), summary=personal.get("summary", ""), experience=experience, education=education, skills=skills, projects=projects)

    def _save_html(self, html_content, path):
        Path(path).parent.mkdir(parents=True, exist_ok=True)
        with open(path, 'w', encoding='utf-8') as f:
            f.write(html_content)

    def _generate_pdf(self, html_content, output_path):
        Path(output_path).parent.mkdir(parents=True, exist_ok=True)
        try:
            self._generate_pdf_weasyprint(html_content, output_path)
            logger.info(f"PDF generated with WeasyPrint: {output_path}")
            return
        except Exception as e:
            logger.warning(f"WeasyPrint PDF generation failed: {e}")
        try:
            self._generate_pdf_selenium(html_content, output_path)
            logger.info(f"PDF generated with Selenium: {output_path}")
            return
        except Exception as e:
            logger.warning(f"Selenium PDF generation failed: {e}")
        try:
            self._generate_pdf_reportlab(html_content, output_path)
            logger.info(f"PDF generated with ReportLab: {output_path}")
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
        story = [Paragraph("Resume PDF - See HTML version for full formatting", styles['Normal'])]
        doc.build(story)
