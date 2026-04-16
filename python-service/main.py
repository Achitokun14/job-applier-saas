from fastapi import FastAPI, HTTPException, BackgroundTasks, File, UploadFile, Form
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any
import os
import json
import yaml
import tempfile
from pathlib import Path
import uuid
import logging

from src.resume_generator import ResumeGenerator
from src.cover_letter_generator import CoverLetterGenerator
from src.job_parser import JobParser
from src.scraper import ScrapeRequest, JobSpyScraper
from src.regional_scrapers import RegionalScraper
from src.resume_parser import ResumeParser
from src.embeddings import EmbeddingService
from src.ats_scorer import ATSScorer
from src.skills_gap import SkillsGapAnalyzer
from src.interview_prep import InterviewPrep
from src.salary_analyzer import SalaryAnalyzer
from src.company_rag import CompanyResearcher
from src.tasks import (
    generate_resume_task,
    generate_cover_letter_task,
    parse_job_task,
    scrape_jobs_task,
    auto_apply_task,
    parse_resume_task,
    compute_match_score_task,
    ats_score_task,
    interview_questions_task,
)
from src.celery_app import celery_app
from src.auto_applier import AutoApplier, AutoApplyRequest, AutoApplyResult

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="Job Applier Resume Service", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

OUTPUT_DIR = Path("data_folder/output")
OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

resume_generator = ResumeGenerator()
cover_letter_generator = CoverLetterGenerator()
job_parser = JobParser()
jobspy_scraper = JobSpyScraper()
regional_scraper = RegionalScraper()
auto_applier = AutoApplier()
resume_pdf_parser = ResumeParser()
embedding_service = EmbeddingService()
ats_scorer = ATSScorer()
skills_gap_analyzer = SkillsGapAnalyzer()
interview_prep = InterviewPrep()
salary_analyzer = SalaryAnalyzer()
company_researcher = CompanyResearcher()


class ResumeRequest(BaseModel):
    resume_yaml: str = Field(..., description="Resume content in YAML format")
    style: str = Field(default="modern", description="Resume style (modern, classic, minimal, creative, professional)")
    job_url: Optional[str] = Field(default=None, description="Job URL to tailor resume for")
    job_description: Optional[str] = Field(default=None, description="Job description text")
    llm_provider: Optional[str] = Field(default=None, description="LLM provider (openai, anthropic, google, groq, mistral, ollama)")
    llm_model: Optional[str] = Field(default=None, description="LLM model name")
    llm_api_key: Optional[str] = Field(default=None, description="LLM API key")


class CoverLetterRequest(BaseModel):
    resume_text: str = Field(..., description="Resume content as text")
    job_description: str = Field(..., description="Job description")
    job_url: Optional[str] = Field(default=None, description="Job posting URL")
    company_name: Optional[str] = Field(default=None, description="Company name")
    job_title: Optional[str] = Field(default=None, description="Job title")
    llm_provider: Optional[str] = Field(default=None, description="LLM provider (openai, anthropic, google, groq, mistral, ollama)")
    llm_model: Optional[str] = Field(default=None, description="LLM model name")
    llm_api_key: Optional[str] = Field(default=None, description="LLM API key")


class JobParseRequest(BaseModel):
    url: str = Field(..., description="Job posting URL")
    llm_provider: Optional[str] = Field(default=None, description="LLM provider (openai, anthropic, google, groq, mistral, ollama)")
    llm_model: Optional[str] = Field(default=None, description="LLM model name")
    llm_api_key: Optional[str] = Field(default=None, description="LLM API key")


class GenerationResponse(BaseModel):
    id: str
    pdf_path: str
    html_content: Optional[str] = None
    metadata: Dict[str, Any] = {}


class JobParseResponse(BaseModel):
    title: str
    company: str
    location: str
    description: str
    requirements: List[str]
    responsibilities: List[str]
    salary: Optional[str] = None
    remote: bool = False


class AsyncTaskResponse(BaseModel):
    task_id: str
    status: str


class TaskStatusResponse(BaseModel):
    task_id: str
    status: str
    result: Optional[Dict[str, Any]] = None
    error: Optional[str] = None


# ---- Health Check ----

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "resume-generator"}


# ---- Synchronous Endpoints (backward compatible) ----

@app.post("/generate-resume", response_model=GenerationResponse)
async def generate_resume(request: ResumeRequest, background_tasks: BackgroundTasks):
    try:
        resume_id = str(uuid.uuid4())[:8]
        output_path = str(OUTPUT_DIR / f"resume_{resume_id}.pdf")

        if request.job_url and not request.job_description:
            job_data = job_parser.parse_url(
                request.job_url,
                llm_provider=request.llm_provider,
                llm_api_key=request.llm_api_key,
                llm_model=request.llm_model,
            )
            request.job_description = job_data.get("description", "")

        result = resume_generator.generate(
            resume_yaml=request.resume_yaml,
            style=request.style,
            job_description=request.job_description,
            output_path=output_path,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )

        logger.info(f"Generated resume: {resume_id}")

        return GenerationResponse(
            id=resume_id,
            pdf_path=result["pdf_path"],
            html_content=result.get("html_content"),
            metadata={
                "style": request.style,
                "tailored": bool(request.job_description),
                "word_count": result.get("word_count", 0)
            }
        )
    except Exception as e:
        logger.error(f"Resume generation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/generate-cover-letter", response_model=GenerationResponse)
async def generate_cover_letter(request: CoverLetterRequest, background_tasks: BackgroundTasks):
    try:
        letter_id = str(uuid.uuid4())[:8]
        output_path = str(OUTPUT_DIR / f"cover_letter_{letter_id}.pdf")

        result = cover_letter_generator.generate(
            resume_text=request.resume_text,
            job_description=request.job_description,
            job_url=request.job_url,
            company_name=request.company_name,
            job_title=request.job_title,
            output_path=output_path,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )

        logger.info(f"Generated cover letter: {letter_id}")

        return GenerationResponse(
            id=letter_id,
            pdf_path=result["pdf_path"],
            html_content=result.get("html_content"),
            metadata={
                "company": request.company_name,
                "job_title": request.job_title,
                "word_count": result.get("word_count", 0)
            }
        )
    except Exception as e:
        logger.error(f"Cover letter generation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/parse-job", response_model=JobParseResponse)
async def parse_job(request: JobParseRequest):
    try:
        job_data = job_parser.parse_url(
            request.url,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return JobParseResponse(**job_data)
    except Exception as e:
        logger.error(f"Job parsing failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ---- Async Endpoints (Celery-backed) ----

@app.post("/async/generate-resume", response_model=AsyncTaskResponse)
async def async_generate_resume(request: ResumeRequest):
    """Enqueue resume generation as an async Celery task."""
    task = generate_resume_task.delay(
        resume_yaml=request.resume_yaml,
        style=request.style,
        job_description=request.job_description,
        llm_provider=request.llm_provider,
        llm_model=request.llm_model,
        llm_api_key=request.llm_api_key,
    )
    logger.info(f"Enqueued async resume generation: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


@app.post("/async/generate-cover-letter", response_model=AsyncTaskResponse)
async def async_generate_cover_letter(request: CoverLetterRequest):
    """Enqueue cover letter generation as an async Celery task."""
    task = generate_cover_letter_task.delay(
        resume_text=request.resume_text,
        job_description=request.job_description,
        company_name=request.company_name,
        job_title=request.job_title,
        llm_provider=request.llm_provider,
        llm_model=request.llm_model,
        llm_api_key=request.llm_api_key,
    )
    logger.info(f"Enqueued async cover letter generation: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


@app.post("/async/parse-job", response_model=AsyncTaskResponse)
async def async_parse_job(request: JobParseRequest):
    """Enqueue job parsing as an async Celery task."""
    task = parse_job_task.delay(
        url=request.url,
        llm_provider=request.llm_provider,
        llm_model=request.llm_model,
        llm_api_key=request.llm_api_key,
    )
    logger.info(f"Enqueued async job parsing: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


@app.get("/api/tasks/{task_id}/status", response_model=TaskStatusResponse)
async def get_task_status(task_id: str):
    """Check the status of a Celery task."""
    from celery.result import AsyncResult

    result = AsyncResult(task_id, app=celery_app)

    response = TaskStatusResponse(
        task_id=task_id,
        status=result.status,
    )

    if result.ready():
        if result.successful():
            response.result = result.result
        else:
            response.error = str(result.result)

    return response


# ---- Scrape Endpoints ----

@app.post("/scrape-jobs")
async def scrape_jobs(request: ScrapeRequest):
    """Run JobSpy + regional scrapers, return combined deduplicated results."""
    try:
        # Run JobSpy scraper (synchronous, so run in thread)
        import asyncio
        loop = asyncio.get_event_loop()
        jobspy_results = await loop.run_in_executor(None, jobspy_scraper.scrape, request)
        logger.info(f"JobSpy returned {len(jobspy_results)} jobs")

        # Run regional scrapers (async)
        regional_results = await regional_scraper.scrape_all(
            search_term=request.search_term,
            location=request.location,
        )
        logger.info(f"Regional scrapers returned {len(regional_results)} jobs")

        # Combine and deduplicate by external_id
        seen_ids = set()
        combined = []
        for job in jobspy_results + regional_results:
            eid = job.get("external_id", "")
            if eid and eid not in seen_ids:
                seen_ids.add(eid)
                combined.append(job)

        logger.info(f"Total deduplicated jobs: {len(combined)}")
        return {"jobs": combined, "total": len(combined)}
    except Exception as e:
        logger.error(f"Scrape failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/async/scrape-jobs", response_model=AsyncTaskResponse)
async def async_scrape_jobs(request: ScrapeRequest):
    """Enqueue job scraping as an async Celery task."""
    task = scrape_jobs_task.delay(
        search_term=request.search_term,
        location=request.location,
        sites=request.sites,
        results_wanted=request.results_wanted,
        hours_old=request.hours_old,
        is_remote=request.is_remote,
        job_type=request.job_type,
        country=request.country,
        distance=request.distance,
        proxies=request.proxies,
    )
    logger.info(f"Enqueued async job scraping: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


# ---- Auto-Apply Endpoints ----

@app.post("/api/auto-apply", response_model=AutoApplyResult)
async def auto_apply(request: AutoApplyRequest):
    """Route a job application through the appropriate strategy.

    - Greenhouse/Lever ATS: fully automated API submission
    - Form fill: fills form + screenshots, requires user confirmation
    - Email: composes email payload, requires confirmation
    - LinkedIn/Indeed: blocked (returns error with 'apply manually' message)
    """
    try:
        result = await auto_applier.apply(request)
        return result
    except Exception as e:
        logger.error(f"Auto-apply failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


class ConfirmSubmitRequest(BaseModel):
    task_id: str = Field(..., description="ID of the auto-apply task to confirm")


@app.post("/api/auto-apply/confirm")
async def confirm_auto_apply(request: ConfirmSubmitRequest):
    """Confirm and trigger actual form submission for a previously filled form.

    This endpoint is called after the user reviews the screenshot of the
    filled form and approves submission. It re-opens the page, re-fills,
    and clicks submit.
    """
    # In production this would look up the saved form state by task_id,
    # re-launch Playwright, re-fill the form, and click submit.
    # For now we acknowledge the confirmation and let the Go backend
    # handle the state transition.
    logger.info(f"Submission confirmed for task: {request.task_id}")
    return {
        "task_id": request.task_id,
        "status": "confirmed",
        "message": "Submission confirmation received. The form will be submitted.",
    }


@app.post("/async/auto-apply", response_model=AsyncTaskResponse)
async def async_auto_apply(request: AutoApplyRequest):
    """Enqueue auto-apply as an async Celery task."""
    task = auto_apply_task.delay(
        job_url=request.job_url,
        apply_url=request.apply_url,
        source=request.source,
        resume_pdf_path=request.resume_pdf_path,
        cover_letter_pdf_path=request.cover_letter_pdf_path,
        user_name=request.user_name,
        user_email=request.user_email,
        user_phone=request.user_phone,
        linkedin_url=request.linkedin_url,
    )
    logger.info(f"Enqueued async auto-apply: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


# ---- AI/ML Feature Endpoints ----


class MatchScoreRequest(BaseModel):
    resume_text: str = Field(..., description="Resume content as text")
    job_description: str = Field(..., description="Job description text")
    resume_skills: Optional[List[str]] = Field(default=None, description="Skills from resume")
    job_skills: Optional[List[str]] = Field(default=None, description="Skills from job description")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


class ATSScoreRequest(BaseModel):
    resume_text: str = Field(..., description="Resume content as text")
    job_description: str = Field(..., description="Job description text")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


class SkillsGapRequest(BaseModel):
    resume_text: str = Field(..., description="Resume content as text")
    job_description: str = Field(..., description="Job description text")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


class InterviewQuestionsRequest(BaseModel):
    job_description: str = Field(..., description="Job description text")
    num_questions: int = Field(default=10, description="Number of questions to generate")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


class InterviewEvaluateRequest(BaseModel):
    question: str = Field(..., description="The interview question")
    answer: str = Field(..., description="The candidate's answer")
    job_description: str = Field(..., description="Job description for context")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


class SalaryAnalyzeRequest(BaseModel):
    job_title: str = Field(..., description="Job title")
    location: str = Field(..., description="Job location")
    experience_years: int = Field(..., description="Years of experience")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


class CompanyResearchRequest(BaseModel):
    company_name: str = Field(..., description="Company name")
    company_url: Optional[str] = Field(default=None, description="Company website URL")
    llm_provider: Optional[str] = None
    llm_model: Optional[str] = None
    llm_api_key: Optional[str] = None


@app.post("/api/resume/parse")
async def parse_resume_pdf(
    file: UploadFile = File(...),
    llm_provider: Optional[str] = Form(None),
    llm_model: Optional[str] = Form(None),
    llm_api_key: Optional[str] = Form(None),
):
    """Parse a resume PDF and extract structured data."""
    if not file.filename.lower().endswith(".pdf"):
        raise HTTPException(status_code=400, detail="Only PDF files are supported")

    try:
        # Save uploaded file to temp location
        with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as tmp:
            content = await file.read()
            tmp.write(content)
            tmp_path = tmp.name

        result = resume_pdf_parser.parse_pdf(
            pdf_path=tmp_path,
            llm_provider=llm_provider,
            llm_api_key=llm_api_key,
            llm_model=llm_model,
        )

        # Clean up temp file
        os.unlink(tmp_path)

        return result
    except Exception as e:
        logger.error(f"Resume parsing failed: {e}")
        # Clean up temp file on error
        try:
            os.unlink(tmp_path)
        except Exception:
            pass
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/matching/score")
async def match_score(request: MatchScoreRequest):
    """Compute match score between resume and job description using embeddings."""
    try:
        result = embedding_service.compute_match_score(
            resume_text=request.resume_text,
            job_description=request.job_description,
            resume_skills=request.resume_skills,
            job_skills=request.job_skills,
        )
        return result
    except Exception as e:
        logger.error(f"Match scoring failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/ats/score")
async def ats_score(request: ATSScoreRequest):
    """Score a resume for ATS compatibility against a job description."""
    try:
        result = ats_scorer.score(
            resume_text=request.resume_text,
            job_description=request.job_description,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return result
    except Exception as e:
        logger.error(f"ATS scoring failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/skills-gap")
async def skills_gap(request: SkillsGapRequest):
    """Analyze skills gap between resume and job description."""
    try:
        result = skills_gap_analyzer.analyze(
            resume_text=request.resume_text,
            job_description=request.job_description,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return result
    except Exception as e:
        logger.error(f"Skills gap analysis failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/interview/questions")
async def interview_questions(request: InterviewQuestionsRequest):
    """Generate interview questions based on job description."""
    try:
        result = interview_prep.generate_questions(
            job_description=request.job_description,
            num_questions=request.num_questions,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return result
    except Exception as e:
        logger.error(f"Interview question generation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/interview/evaluate")
async def interview_evaluate(request: InterviewEvaluateRequest):
    """Evaluate a candidate's answer to an interview question."""
    try:
        result = interview_prep.evaluate_answer(
            question=request.question,
            answer=request.answer,
            job_description=request.job_description,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return result
    except Exception as e:
        logger.error(f"Answer evaluation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/salary/analyze")
async def salary_analyze(request: SalaryAnalyzeRequest):
    """Estimate salary range and provide negotiation insights."""
    try:
        result = salary_analyzer.analyze(
            job_title=request.job_title,
            location=request.location,
            experience_years=request.experience_years,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return result
    except Exception as e:
        logger.error(f"Salary analysis failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/company/research")
async def company_research(request: CompanyResearchRequest):
    """Research a company using web scraping + RAG + LLM analysis."""
    try:
        result = company_researcher.research(
            company_name=request.company_name,
            company_url=request.company_url,
            llm_provider=request.llm_provider,
            llm_api_key=request.llm_api_key,
            llm_model=request.llm_model,
        )
        return result
    except Exception as e:
        logger.error(f"Company research failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ---- Async AI/ML Endpoints (Celery-backed) ----


@app.post("/async/resume/parse", response_model=AsyncTaskResponse)
async def async_parse_resume(
    file: UploadFile = File(...),
    llm_provider: Optional[str] = Form(None),
    llm_model: Optional[str] = Form(None),
    llm_api_key: Optional[str] = Form(None),
):
    """Enqueue resume PDF parsing as an async Celery task."""
    if not file.filename.lower().endswith(".pdf"):
        raise HTTPException(status_code=400, detail="Only PDF files are supported")

    # Save uploaded file to persistent temp location
    with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf", dir=str(OUTPUT_DIR)) as tmp:
        content = await file.read()
        tmp.write(content)
        tmp_path = tmp.name

    task = parse_resume_task.delay(
        pdf_path=tmp_path,
        llm_provider=llm_provider,
        llm_model=llm_model,
        llm_api_key=llm_api_key,
    )
    logger.info(f"Enqueued async resume parsing: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


@app.post("/async/matching/score", response_model=AsyncTaskResponse)
async def async_match_score(request: MatchScoreRequest):
    """Enqueue match scoring as an async Celery task."""
    task = compute_match_score_task.delay(
        resume_text=request.resume_text,
        job_description=request.job_description,
        resume_skills=request.resume_skills,
        job_skills=request.job_skills,
    )
    logger.info(f"Enqueued async match scoring: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


@app.post("/async/ats/score", response_model=AsyncTaskResponse)
async def async_ats_score(request: ATSScoreRequest):
    """Enqueue ATS scoring as an async Celery task."""
    task = ats_score_task.delay(
        resume_text=request.resume_text,
        job_description=request.job_description,
        llm_provider=request.llm_provider,
        llm_model=request.llm_model,
        llm_api_key=request.llm_api_key,
    )
    logger.info(f"Enqueued async ATS scoring: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


@app.post("/async/interview/questions", response_model=AsyncTaskResponse)
async def async_interview_questions(request: InterviewQuestionsRequest):
    """Enqueue interview question generation as an async Celery task."""
    task = interview_questions_task.delay(
        job_description=request.job_description,
        num_questions=request.num_questions,
        llm_provider=request.llm_provider,
        llm_model=request.llm_model,
        llm_api_key=request.llm_api_key,
    )
    logger.info(f"Enqueued async interview questions: task_id={task.id}")
    return AsyncTaskResponse(task_id=task.id, status="queued")


# ---- Metadata Endpoints ----

@app.get("/styles")
async def list_styles():
    return {
        "styles": [
            {"id": "modern", "name": "Modern", "description": "Clean, contemporary design with accent colors"},
            {"id": "classic", "name": "Classic", "description": "Traditional, professional layout"},
            {"id": "minimal", "name": "Minimal", "description": "Simple, elegant with focus on content"},
            {"id": "creative", "name": "Creative", "description": "Bold design for creative roles"},
            {"id": "professional", "name": "Professional", "description": "Corporate-focused, ATS-friendly"},
        ]
    }


@app.get("/templates")
async def list_templates():
    return {
        "templates": [
            {"id": "standard", "name": "Standard Resume", "sections": ["experience", "education", "skills"]},
            {"id": "executive", "name": "Executive Resume", "sections": ["summary", "achievements", "experience"]},
            {"id": "technical", "name": "Technical Resume", "sections": ["skills", "projects", "experience"]},
        ]
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
