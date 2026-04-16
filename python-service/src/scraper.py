import hashlib
import logging
from typing import Optional
from pydantic import BaseModel, Field

logger = logging.getLogger(__name__)


class ScrapeRequest(BaseModel):
    search_term: str
    location: Optional[str] = None
    sites: list[str] = Field(default=["indeed", "linkedin", "glassdoor", "google"])
    results_wanted: int = 50
    hours_old: int = 72
    is_remote: Optional[bool] = None
    job_type: Optional[str] = None
    country: Optional[str] = None
    distance: Optional[int] = 50
    proxies: Optional[list[str]] = None


class ScrapedJob(BaseModel):
    external_id: str
    source: str
    title: str
    company: str
    location: str
    description: str
    url: str
    remote: bool = False
    salary: str = ""
    posted_at: str = ""


class JobSpyScraper:
    def scrape(self, request: ScrapeRequest) -> list[dict]:
        from jobspy import scrape_jobs

        try:
            kwargs = {
                "site_name": request.sites,
                "search_term": request.search_term,
                "results_wanted": request.results_wanted,
                "hours_old": request.hours_old,
            }
            if request.location:
                kwargs["location"] = request.location
            if request.is_remote is not None:
                kwargs["is_remote"] = request.is_remote
            if request.job_type:
                kwargs["job_type"] = request.job_type
            if request.country:
                kwargs["country_indeed"] = request.country
            if request.distance:
                kwargs["distance"] = request.distance
            if request.proxies:
                kwargs["proxies"] = request.proxies

            df = scrape_jobs(**kwargs)

            jobs = []
            for _, row in df.iterrows():
                external_id = self._generate_id(row)
                job = {
                    "external_id": external_id,
                    "source": str(row.get("site", "unknown")),
                    "title": str(row.get("title", "")),
                    "company": str(row.get("company", "")),
                    "location": str(row.get("location", "")),
                    "description": str(row.get("description", ""))[:5000],
                    "url": str(row.get("job_url", "")),
                    "remote": bool(row.get("is_remote", False)),
                    "salary": self._format_salary(row),
                    "posted_at": str(row.get("date_posted", "")),
                }
                jobs.append(job)
            return jobs
        except Exception as e:
            logger.error(f"JobSpy scrape failed: {e}")
            return []

    def _generate_id(self, row) -> str:
        raw = f"{row.get('site', '')}-{row.get('id', row.get('job_url', ''))}"
        return hashlib.sha256(raw.encode()).hexdigest()[:16]

    def _format_salary(self, row) -> str:
        min_s = row.get("min_amount")
        max_s = row.get("max_amount")
        if min_s and max_s and not (str(min_s) == "nan" or str(max_s) == "nan"):
            try:
                return f"${int(float(min_s)):,}-${int(float(max_s)):,}"
            except (ValueError, TypeError):
                return ""
        return ""
