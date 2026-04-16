import os
import re
from typing import Dict, Any, List, Optional

from pydantic import BaseModel, Field


class ParsedJob(BaseModel):
    title: str = Field(description="Job title")
    company: str = Field(description="Company name")
    location: str = Field(description="Job location")
    description: str = Field(description="Full job description")
    requirements: List[str] = Field(default_factory=list, description="List of requirements")
    responsibilities: List[str] = Field(default_factory=list, description="List of responsibilities")
    salary: Optional[str] = Field(default=None, description="Salary range if mentioned")
    remote: bool = Field(default=False, description="Whether remote work is offered")


class JobParser:
    def __init__(self):
        self.llm_api_key = os.getenv("LLM_API_KEY", "")
        self.llm_model = os.getenv("LLM_MODEL", "gpt-4o-mini")
        self.llm_provider = os.getenv("LLM_PROVIDER", "openai")

    def parse_url(
        self,
        url: str,
        llm_provider: Optional[str] = None,
        llm_api_key: Optional[str] = None,
        llm_model: Optional[str] = None,
    ) -> Dict[str, Any]:
        api_key = llm_api_key or self.llm_api_key
        if api_key:
            return self._parse_with_llm(
                url,
                provider=llm_provider or self.llm_provider,
                api_key=api_key,
                model=llm_model or self.llm_model,
            )
        else:
            return self._parse_basic(url)

    def _parse_with_llm(
        self,
        url: str,
        provider: str,
        api_key: str,
        model: str,
    ) -> Dict[str, Any]:
        try:
            from selenium import webdriver
            from selenium.webdriver.chrome.options import Options
            from selenium.webdriver.common.by import By
            from selenium.webdriver.support.ui import WebDriverWait
            from selenium.webdriver.support import expected_conditions as EC
            from src.llm_client import get_client

            chrome_options = Options()
            chrome_options.add_argument('--headless')
            chrome_options.add_argument('--no-sandbox')
            chrome_options.add_argument('--disable-dev-shm-usage')

            driver = webdriver.Chrome(options=chrome_options)
            driver.get(url)

            wait = WebDriverWait(driver, 10)
            wait.until(EC.presence_of_element_located((By.TAG_NAME, "body")))

            page_text = driver.find_element(By.TAG_NAME, "body").text
            driver.quit()

            client, full_model = get_client(
                provider=provider,
                api_key=api_key,
                model=model,
            )

            parsed_job = client.chat.completions.create(
                model=full_model,
                response_model=ParsedJob,
                temperature=0,
                messages=[
                    {
                        "role": "user",
                        "content": (
                            "Parse this job posting and extract structured information.\n\n"
                            f"{page_text[:5000]}"
                        ),
                    }
                ],
            )

            return parsed_job.model_dump()

        except Exception as e:
            print(f"LLM parsing failed: {e}")
            return self._parse_basic(url)

    def _parse_basic(self, url: str) -> Dict[str, Any]:
        try:
            from selenium import webdriver
            from selenium.webdriver.chrome.options import Options
            from selenium.webdriver.common.by import By

            chrome_options = Options()
            chrome_options.add_argument('--headless')
            chrome_options.add_argument('--no-sandbox')
            chrome_options.add_argument('--disable-dev-shm-usage')

            driver = webdriver.Chrome(options=chrome_options)
            driver.get(url)

            page_text = driver.find_element(By.TAG_NAME, "body").text
            driver.quit()

            return self._extract_from_text(page_text, url)

        except Exception as e:
            print(f"Basic parsing failed: {e}")
            return {
                "title": "Unknown Position",
                "company": "Unknown Company",
                "location": "Unknown Location",
                "description": f"Job posting from {url}",
                "requirements": [],
                "responsibilities": [],
                "salary": None,
                "remote": False
            }

    def _extract_from_text(self, text: str, url: str) -> Dict[str, Any]:
        title_patterns = [
            r'(?:job\s*title|position)[:\s]+(.+?)(?:\n|$)',
            r'^([A-Z][A-Za-z\s]+(?:Engineer|Developer|Manager|Designer|Analyst))',
        ]

        title = "Unknown Position"
        for pattern in title_patterns:
            match = re.search(pattern, text, re.IGNORECASE | re.MULTILINE)
            if match:
                title = match.group(1).strip()
                break

        company = "Unknown Company"
        if "linkedin.com" in url:
            company_match = re.search(r'at\s+([A-Z][A-Za-z\s&]+)', text)
            if company_match:
                company = company_match.group(1).strip()

        remote = bool(re.search(r'\b(remote|work\s*from\s*home|wfh)\b', text, re.IGNORECASE))

        requirements = re.findall(
            r'(?:requirements?|qualifications?|must\s*have)[:\s]*\n((?:[-•*]\s*.+\n?)+)',
            text, re.IGNORECASE
        )

        req_list = []
        if requirements:
            req_list = [r.strip() for r in re.split(r'[-•*]\s*', requirements[0]) if r.strip()]

        return {
            "title": title,
            "company": company,
            "location": "See posting for details",
            "description": text[:2000],
            "requirements": req_list[:10],
            "responsibilities": [],
            "salary": None,
            "remote": remote
        }
