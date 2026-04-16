import logging
from pydantic import BaseModel, Field
from typing import Optional

logger = logging.getLogger(__name__)


class MissingSkill(BaseModel):
    name: str = Field(description="Name of the missing skill")
    importance: str = Field(
        description="How important this skill is: 'required' or 'preferred'"
    )
    recommendation: str = Field(
        description="Actionable recommendation for acquiring this skill"
    )


class PartialMatch(BaseModel):
    resume_skill: str = Field(description="Skill found in the resume")
    job_skill: str = Field(description="Related skill from the job description")
    match_level: str = Field(
        description="How close the match is: 'strong', 'moderate', or 'weak'"
    )


class SkillsGapResult(BaseModel):
    present_skills: list[str] = Field(
        default_factory=list,
        description="Skills from the job description that are present in the resume",
    )
    missing_skills: list[MissingSkill] = Field(
        default_factory=list,
        description="Skills required by the job but missing from the resume",
    )
    partial_matches: list[PartialMatch] = Field(
        default_factory=list,
        description="Skills that partially match between resume and job",
    )


class SkillsGapAnalyzer:
    def analyze(
        self,
        resume_text: str,
        job_description: str,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Analyze the skills gap between a resume and a job description."""
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )

            result = client.chat.completions.create(
                model=model,
                response_model=SkillsGapResult,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are a career skills analyst. Compare the candidate's resume "
                            "against the job description to identify:\n"
                            "1. Skills present in both the resume and job description\n"
                            "2. Skills required by the job but missing from the resume, "
                            "with importance level (required/preferred) and actionable recommendations\n"
                            "3. Partial matches where the candidate has a related but not exact skill\n\n"
                            "Be thorough and consider both technical and soft skills. "
                            "Include specific tools, technologies, frameworks, and methodologies."
                        ),
                    },
                    {
                        "role": "user",
                        "content": (
                            f"Resume:\n{resume_text[:4000]}\n\n"
                            f"Job Description:\n{job_description[:3000]}"
                        ),
                    },
                ],
            )

            return {
                "present_skills": result.present_skills,
                "missing_skills": [s.model_dump() for s in result.missing_skills],
                "partial_matches": [m.model_dump() for m in result.partial_matches],
                "summary": {
                    "total_required": len(result.present_skills) + len(
                        [s for s in result.missing_skills if s.importance == "required"]
                    ),
                    "matched": len(result.present_skills),
                    "missing_required": len(
                        [s for s in result.missing_skills if s.importance == "required"]
                    ),
                    "missing_preferred": len(
                        [s for s in result.missing_skills if s.importance == "preferred"]
                    ),
                    "partial": len(result.partial_matches),
                },
            }

        except Exception as e:
            logger.error(f"Skills gap analysis failed: {e}")
            raise
