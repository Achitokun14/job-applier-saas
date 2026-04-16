import logging
import re
from pydantic import BaseModel, Field
from typing import Optional

logger = logging.getLogger(__name__)


class ATSSuggestions(BaseModel):
    """LLM-generated ATS improvement suggestions."""

    suggestions: list[str] = Field(
        default_factory=list,
        description="Specific, actionable suggestions to improve ATS score",
    )
    missing_keywords: list[str] = Field(
        default_factory=list,
        description="Important keywords from the job description missing in the resume",
    )
    format_issues: list[str] = Field(
        default_factory=list,
        description="Formatting issues that may affect ATS parsing",
    )


class ATSScorer:
    # Common action verbs used in resumes
    ACTION_VERBS = {
        "managed", "developed", "led", "built", "designed", "implemented",
        "created", "improved", "increased", "reduced", "delivered",
        "launched", "optimized", "established", "coordinated", "executed",
        "directed", "achieved", "generated", "maintained", "organized",
        "produced", "resolved", "streamlined", "supervised", "trained",
        "analyzed", "collaborated", "contributed", "engineered", "facilitated",
        "initiated", "mentored", "negotiated", "oversaw", "pioneered",
        "spearheaded", "transformed", "automated", "architected", "scaled",
    }

    # Standard resume section headings
    SECTION_HEADINGS = {
        "experience", "work experience", "professional experience",
        "education", "skills", "technical skills", "summary",
        "professional summary", "certifications", "projects",
        "achievements", "awards", "publications", "languages",
        "objective", "career objective", "qualifications",
    }

    def score(
        self,
        resume_text: str,
        job_description: str,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Score a resume against a job description for ATS compatibility."""
        criteria = []
        resume_lower = resume_text.lower()
        job_lower = job_description.lower()

        # 1. Keyword Match (30%)
        job_words = set(re.findall(r"\b[a-zA-Z]{3,}\b", job_lower))
        # Filter out very common English words
        stopwords = {
            "the", "and", "for", "are", "but", "not", "you", "all",
            "can", "had", "her", "was", "one", "our", "out", "has",
            "have", "will", "with", "this", "that", "from", "they",
            "been", "said", "each", "which", "their", "about", "would",
            "make", "like", "just", "over", "such", "take", "other",
            "than", "then", "them", "some", "could", "into", "year",
            "also", "back", "after", "work", "first", "well", "even",
            "where", "what", "when", "who", "how", "more", "should",
        }
        job_keywords = job_words - stopwords
        if job_keywords:
            resume_words = set(re.findall(r"\b[a-zA-Z]{3,}\b", resume_lower))
            matched = job_keywords & resume_words
            keyword_score = min(len(matched) / len(job_keywords) * 100, 100)
        else:
            keyword_score = 50
        criteria.append({
            "name": "keyword_match",
            "score": round(keyword_score),
            "weight": 30,
            "description": "Job description keywords found in resume",
        })

        # 2. Format Compliance (20%)
        found_sections = 0
        for heading in self.SECTION_HEADINGS:
            # Check for heading as a standalone line or with common formatting
            if re.search(
                rf"(?:^|\n)\s*(?:#+\s*)?{re.escape(heading)}\s*(?:\n|:|\|)",
                resume_lower,
            ):
                found_sections += 1
        format_score = min(found_sections / 4 * 100, 100)  # Expect at least 4 sections
        criteria.append({
            "name": "format_compliance",
            "score": round(format_score),
            "weight": 20,
            "description": "Standard section headings present",
        })

        # 3. Content Density (15%)
        word_count = len(resume_text.split())
        if 300 <= word_count <= 800:
            density_score = 100
        elif word_count < 300:
            density_score = max(word_count / 300 * 100, 10)
        else:
            # Penalize slightly for being too long, but not too harshly
            density_score = max(100 - (word_count - 800) / 10, 40)
        criteria.append({
            "name": "content_density",
            "score": round(density_score),
            "weight": 15,
            "description": f"Word count: {word_count} (optimal: 300-800)",
        })

        # 4. Action Verbs (10%)
        lines = resume_text.split("\n")
        bullet_lines = [
            line.strip()
            for line in lines
            if line.strip().startswith(("-", "*", "•")) or re.match(r"^\d+\.", line.strip())
        ]
        if bullet_lines:
            action_count = sum(
                1
                for line in bullet_lines
                if any(
                    line.lower().lstrip("-*•0123456789. ").startswith(verb)
                    for verb in self.ACTION_VERBS
                )
            )
            action_score = min(action_count / max(len(bullet_lines), 1) * 100, 100)
        else:
            action_score = 20  # No bullet points found
        criteria.append({
            "name": "action_verbs",
            "score": round(action_score),
            "weight": 10,
            "description": "Bullet points starting with action verbs",
        })

        # 5. Quantification (10%)
        numbers = re.findall(r"\b\d+[%+]?\b", resume_text)
        percentages = re.findall(r"\d+%", resume_text)
        dollar_amounts = re.findall(r"\$[\d,]+", resume_text)
        quant_count = len(numbers) + len(percentages) + len(dollar_amounts)
        quant_score = min(quant_count / 5 * 100, 100)  # Expect at least 5 quantified items
        criteria.append({
            "name": "quantification",
            "score": round(quant_score),
            "weight": 10,
            "description": f"Quantified achievements: {quant_count} numbers/percentages found",
        })

        # 6. Skills Section (15%)
        has_skills_section = bool(
            re.search(
                r"(?:^|\n)\s*(?:#+\s*)?(?:skills|technical skills|core competencies)\s*(?:\n|:|\|)",
                resume_lower,
            )
        )
        if has_skills_section and job_keywords:
            # Check how many job keywords appear near the skills section
            skills_match = re.search(
                r"(?:skills|technical skills|core competencies)[\s\S]{0,500}",
                resume_lower,
            )
            if skills_match:
                skills_text = skills_match.group()
                skills_words = set(re.findall(r"\b[a-zA-Z]{3,}\b", skills_text))
                skills_overlap = len(skills_words & job_keywords)
                skills_score = min(skills_overlap / min(len(job_keywords), 15) * 100, 100)
            else:
                skills_score = 40
        elif has_skills_section:
            skills_score = 60
        else:
            skills_score = 10
        criteria.append({
            "name": "skills_section",
            "score": round(skills_score),
            "weight": 15,
            "description": "Skills section with relevant keywords",
        })

        # Calculate overall weighted score
        overall = sum(c["score"] * c["weight"] / 100 for c in criteria)

        # Generate LLM-based suggestions
        suggestions = self._generate_suggestions(
            resume_text, job_description, criteria, overall,
            llm_provider, llm_api_key, llm_model,
        )

        return {
            "overall_score": round(overall),
            "criteria": criteria,
            "suggestions": suggestions,
        }

    def _generate_suggestions(
        self,
        resume_text: str,
        job_description: str,
        criteria: list[dict],
        overall_score: float,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> list[str]:
        """Generate LLM-powered suggestions for improving ATS score."""
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )

            criteria_summary = "\n".join(
                f"- {c['name']}: {c['score']}/100 (weight: {c['weight']}%)"
                for c in criteria
            )

            result = client.chat.completions.create(
                model=model,
                response_model=ATSSuggestions,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are an ATS (Applicant Tracking System) expert. "
                            "Analyze the resume against the job description and provide "
                            "specific, actionable suggestions to improve the ATS score. "
                            "Focus on missing keywords, formatting issues, and content improvements."
                        ),
                    },
                    {
                        "role": "user",
                        "content": (
                            f"Resume (truncated):\n{resume_text[:3000]}\n\n"
                            f"Job Description (truncated):\n{job_description[:2000]}\n\n"
                            f"Current ATS Score: {overall_score}/100\n"
                            f"Criteria Breakdown:\n{criteria_summary}\n\n"
                            "Provide specific suggestions to improve the ATS score."
                        ),
                    },
                ],
            )

            all_suggestions = result.suggestions.copy()
            if result.missing_keywords:
                all_suggestions.append(
                    f"Add these missing keywords: {', '.join(result.missing_keywords[:10])}"
                )
            if result.format_issues:
                all_suggestions.extend(result.format_issues)

            return all_suggestions

        except Exception as e:
            logger.warning(f"LLM suggestion generation failed: {e}")
            # Fall back to rule-based suggestions
            suggestions = []
            for c in criteria:
                if c["score"] < 50:
                    if c["name"] == "keyword_match":
                        suggestions.append(
                            "Add more keywords from the job description to your resume"
                        )
                    elif c["name"] == "format_compliance":
                        suggestions.append(
                            "Add standard section headings: Experience, Education, Skills"
                        )
                    elif c["name"] == "content_density":
                        suggestions.append(
                            "Adjust resume length to 300-800 words for optimal ATS parsing"
                        )
                    elif c["name"] == "action_verbs":
                        suggestions.append(
                            "Start bullet points with action verbs like 'managed', 'developed', 'led'"
                        )
                    elif c["name"] == "quantification":
                        suggestions.append(
                            "Add numbers and percentages to quantify your achievements"
                        )
                    elif c["name"] == "skills_section":
                        suggestions.append(
                            "Add a dedicated Skills section with keywords from the job description"
                        )
            return suggestions
