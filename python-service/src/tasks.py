import asyncio
import os
import uuid
from pathlib import Path
from typing import Optional

from src.celery_app import celery_app
from src.resume_generator import ResumeGenerator
from src.cover_letter_generator import CoverLetterGenerator
from src.job_parser import JobParser
from src.scraper import ScrapeRequest, JobSpyScraper
from src.regional_scrapers import RegionalScraper
from src.auto_applier import AutoApplier, AutoApplyRequest
from src.resume_parser import ResumeParser
from src.embeddings import EmbeddingService
from src.ats_scorer import ATSScorer
from src.interview_prep import InterviewPrep

OUTPUT_DIR = Path("data_folder/output")
OUTPUT_DIR.mkdir(parents=True, exist_ok=True)


@celery_app.task(bind=True, name="generate_resume_task", max_retries=2)
def generate_resume_task(
    self,
    resume_yaml: str,
    style: str = "modern",
    job_description: str = None,
    llm_provider: str = None,
    llm_model: str = None,
    llm_api_key: str = None,
):
    """Generate a resume PDF asynchronously."""
    try:
        resume_id = str(uuid.uuid4())[:8]
        output_path = str(OUTPUT_DIR / f"resume_{resume_id}.pdf")

        generator = ResumeGenerator()
        result = generator.generate(
            resume_yaml=resume_yaml,
            style=style,
            job_description=job_description,
            output_path=output_path,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )

        return {
            "id": resume_id,
            "pdf_path": result["pdf_path"],
            "html_content": result.get("html_content"),
            "word_count": result.get("word_count", 0),
            "style": style,
        }
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)


@celery_app.task(bind=True, name="generate_cover_letter_task", max_retries=2)
def generate_cover_letter_task(
    self,
    resume_text: str,
    job_description: str,
    company_name: str = None,
    job_title: str = None,
    llm_provider: str = None,
    llm_model: str = None,
    llm_api_key: str = None,
):
    """Generate a cover letter PDF asynchronously."""
    try:
        letter_id = str(uuid.uuid4())[:8]
        output_path = str(OUTPUT_DIR / f"cover_letter_{letter_id}.pdf")

        generator = CoverLetterGenerator()
        result = generator.generate(
            resume_text=resume_text,
            job_description=job_description,
            company_name=company_name,
            job_title=job_title,
            output_path=output_path,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )

        return {
            "id": letter_id,
            "pdf_path": result["pdf_path"],
            "html_content": result.get("html_content"),
            "word_count": result.get("word_count", 0),
            "company_name": company_name,
            "job_title": job_title,
        }
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)


@celery_app.task(bind=True, name="parse_job_task", max_retries=2)
def parse_job_task(
    self,
    url: str,
    llm_provider: str = None,
    llm_model: str = None,
    llm_api_key: str = None,
):
    """Parse a job posting URL asynchronously."""
    try:
        parser = JobParser()
        result = parser.parse_url(
            url,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )

        return result
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)


@celery_app.task(bind=True, name="scrape_jobs_task", max_retries=2)
def scrape_jobs_task(
    self,
    search_term: str,
    location: Optional[str] = None,
    sites: Optional[list] = None,
    results_wanted: int = 50,
    hours_old: int = 72,
    is_remote: Optional[bool] = None,
    job_type: Optional[str] = None,
    country: Optional[str] = None,
    distance: Optional[int] = 50,
    proxies: Optional[list] = None,
):
    """Scrape jobs using JobSpy + regional scrapers asynchronously."""
    try:
        if sites is None:
            sites = ["indeed", "linkedin", "glassdoor", "google"]

        request = ScrapeRequest(
            search_term=search_term,
            location=location,
            sites=sites,
            results_wanted=results_wanted,
            hours_old=hours_old,
            is_remote=is_remote,
            job_type=job_type,
            country=country,
            distance=distance,
            proxies=proxies,
        )

        # Run JobSpy scraper (synchronous)
        scraper = JobSpyScraper()
        jobspy_results = scraper.scrape(request)

        # Run regional scrapers (async -- run in new event loop for Celery)
        regional = RegionalScraper()
        loop = asyncio.new_event_loop()
        try:
            regional_results = loop.run_until_complete(
                regional.scrape_all(search_term=search_term, location=location)
            )
        finally:
            loop.close()

        # Combine and deduplicate by external_id
        seen_ids = set()
        combined = []
        for job in jobspy_results + regional_results:
            eid = job.get("external_id", "")
            if eid and eid not in seen_ids:
                seen_ids.add(eid)
                combined.append(job)

        return {"jobs": combined, "total": len(combined)}
    except Exception as exc:
        raise self.retry(exc=exc, countdown=10)


@celery_app.task(bind=True, name="auto_apply_task", max_retries=1)
def auto_apply_task(
    self,
    job_url: str,
    apply_url: str,
    source: str,
    resume_pdf_path: str,
    cover_letter_pdf_path: Optional[str] = None,
    user_name: str = "",
    user_email: str = "",
    user_phone: Optional[str] = None,
    linkedin_url: Optional[str] = None,
):
    """Run auto-apply via the appropriate strategy (Celery background task)."""
    try:
        request = AutoApplyRequest(
            job_url=job_url,
            apply_url=apply_url,
            source=source,
            resume_pdf_path=resume_pdf_path,
            cover_letter_pdf_path=cover_letter_pdf_path,
            user_name=user_name,
            user_email=user_email,
            user_phone=user_phone,
            linkedin_url=linkedin_url,
        )

        applier = AutoApplier()
        loop = asyncio.new_event_loop()
        try:
            result = loop.run_until_complete(applier.apply(request))
        finally:
            loop.close()

        return {
            "success": result.success,
            "method": result.method,
            "confirmation_id": result.confirmation_id,
            "screenshot_path": result.screenshot_path,
            "error": result.error,
            "requires_confirmation": result.requires_confirmation,
        }
    except Exception as exc:
        raise self.retry(exc=exc, countdown=10)


# ---- AI/ML Feature Tasks ----


@celery_app.task(bind=True, name="parse_resume_task", max_retries=2)
def parse_resume_task(
    self,
    pdf_path: str,
    llm_provider: str = None,
    llm_model: str = None,
    llm_api_key: str = None,
):
    """Parse a resume PDF asynchronously."""
    try:
        parser = ResumeParser()
        result = parser.parse_pdf(
            pdf_path=pdf_path,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )

        # Clean up temp file after parsing
        try:
            os.unlink(pdf_path)
        except Exception:
            pass

        return result
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)


@celery_app.task(bind=True, name="compute_match_score_task", max_retries=2)
def compute_match_score_task(
    self,
    resume_text: str,
    job_description: str,
    resume_skills: list = None,
    job_skills: list = None,
):
    """Compute resume-job match score asynchronously."""
    try:
        service = EmbeddingService()
        result = service.compute_match_score(
            resume_text=resume_text,
            job_description=job_description,
            resume_skills=resume_skills,
            job_skills=job_skills,
        )
        return result
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)


@celery_app.task(bind=True, name="ats_score_task", max_retries=2)
def ats_score_task(
    self,
    resume_text: str,
    job_description: str,
    llm_provider: str = None,
    llm_model: str = None,
    llm_api_key: str = None,
):
    """Score resume for ATS compatibility asynchronously."""
    try:
        scorer = ATSScorer()
        result = scorer.score(
            resume_text=resume_text,
            job_description=job_description,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )
        return result
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)


@celery_app.task(bind=True, name="interview_questions_task", max_retries=2)
def interview_questions_task(
    self,
    job_description: str,
    num_questions: int = 10,
    llm_provider: str = None,
    llm_model: str = None,
    llm_api_key: str = None,
):
    """Generate interview questions asynchronously."""
    try:
        prep = InterviewPrep()
        result = prep.generate_questions(
            job_description=job_description,
            num_questions=num_questions,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )
        return result
    except Exception as exc:
        raise self.retry(exc=exc, countdown=5)
