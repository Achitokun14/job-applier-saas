import logging
from pydantic import BaseModel, Field
from typing import List, Optional

logger = logging.getLogger(__name__)


class SalaryEstimate(BaseModel):
    low: int = Field(description="Low end of salary range")
    median: int = Field(description="Median/typical salary")
    high: int = Field(description="High end of salary range")
    currency: str = Field(default="USD", description="Currency code")
    factors: List[str] = Field(
        default_factory=list,
        description="Key factors affecting salary for this role",
    )
    negotiation_tips: List[str] = Field(
        default_factory=list,
        description="Practical tips for salary negotiation",
    )
    market_context: str = Field(
        default="",
        description="Brief context about the current market for this role",
    )


class SalaryAnalyzer:
    def analyze(
        self,
        job_title: str,
        location: str,
        experience_years: int,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Estimate salary range and provide negotiation insights."""
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )

            result = client.chat.completions.create(
                model=model,
                response_model=SalaryEstimate,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are a compensation analyst with deep knowledge of tech and "
                            "professional salary markets. Provide realistic salary estimates "
                            "based on current market data. Consider:\n"
                            "- Location-based cost of living adjustments\n"
                            "- Experience level impact on compensation\n"
                            "- Industry standards and market demand\n"
                            "- Remote vs on-site pay differences\n\n"
                            "Provide salary in annual figures. Be specific about factors "
                            "that affect this particular role's compensation and give "
                            "actionable negotiation advice."
                        ),
                    },
                    {
                        "role": "user",
                        "content": (
                            f"Estimate the salary range for:\n"
                            f"- Job Title: {job_title}\n"
                            f"- Location: {location}\n"
                            f"- Years of Experience: {experience_years}\n\n"
                            "Provide low, median, and high salary estimates with "
                            "factors and negotiation tips."
                        ),
                    },
                ],
            )

            return result.model_dump()

        except Exception as e:
            logger.error(f"Salary analysis failed: {e}")
            raise
