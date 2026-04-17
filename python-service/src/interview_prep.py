import logging
from pydantic import BaseModel, Field
from typing import List, Optional

logger = logging.getLogger(__name__)


# ---- Pydantic models for structured output ----


class InterviewQuestion(BaseModel):
    question: str = Field(description="The interview question")
    category: str = Field(description="Category: behavioral, technical, or situational")
    difficulty: str = Field(description="Difficulty: easy, medium, or hard")
    tips: str = Field(description="Brief tips for answering this question well")


class InterviewQuestions(BaseModel):
    behavioral: List[InterviewQuestion] = Field(
        default_factory=list,
        description="Behavioral interview questions (STAR method)",
    )
    technical: List[InterviewQuestion] = Field(
        default_factory=list,
        description="Technical interview questions about skills and knowledge",
    )
    situational: List[InterviewQuestion] = Field(
        default_factory=list,
        description="Situational/hypothetical scenario questions",
    )


class STARAnalysis(BaseModel):
    situation: str = Field(default="", description="The situation described")
    task: str = Field(default="", description="The task or challenge identified")
    action: str = Field(default="", description="Actions taken by the candidate")
    result: str = Field(default="", description="Results or outcomes achieved")
    completeness: str = Field(
        default="",
        description="Assessment of how complete the STAR response is",
    )


class AnswerEvaluation(BaseModel):
    score: int = Field(
        description="Score from 0-100 for the answer quality", ge=0, le=100
    )
    feedback: str = Field(description="Detailed feedback on the answer")
    improved_answer: str = Field(
        description="An improved version of the answer the candidate could give"
    )
    star_analysis: STARAnalysis = Field(
        default_factory=STARAnalysis,
        description="STAR method analysis of the answer",
    )
    strengths: List[str] = Field(
        default_factory=list,
        description="Strengths of the answer",
    )
    weaknesses: List[str] = Field(
        default_factory=list,
        description="Areas for improvement",
    )


# ---- Service class ----


class InterviewPrep:
    def generate_questions(
        self,
        job_description: str,
        num_questions: int = 10,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Generate categorized interview questions based on a job description."""
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )

            result = client.chat.completions.create(
                model=model,
                response_model=InterviewQuestions,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are an expert interview coach. Generate realistic interview "
                            "questions for the given job description. Create a mix of:\n"
                            "- Behavioral questions (using STAR method expectations)\n"
                            "- Technical questions (testing skills and knowledge)\n"
                            "- Situational questions (hypothetical scenarios)\n\n"
                            "Each question should include its category, difficulty level, "
                            "and brief tips for answering well. Make the questions specific "
                            "to the role and industry."
                        ),
                    },
                    {
                        "role": "user",
                        "content": (
                            f"Generate {num_questions} interview questions for this role:\n\n"
                            f"{job_description[:4000]}"
                        ),
                    },
                ],
            )

            return {
                "behavioral": [q.model_dump() for q in result.behavioral],
                "technical": [q.model_dump() for q in result.technical],
                "situational": [q.model_dump() for q in result.situational],
                "total_questions": (
                    len(result.behavioral)
                    + len(result.technical)
                    + len(result.situational)
                ),
            }

        except Exception as e:
            logger.error(f"Interview question generation failed: {e}")
            raise

    def evaluate_answer(
        self,
        question: str,
        answer: str,
        job_description: str,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Evaluate a candidate's answer to an interview question."""
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )

            result = client.chat.completions.create(
                model=model,
                response_model=AnswerEvaluation,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are an expert interview coach evaluating a candidate's answer. "
                            "Provide:\n"
                            "1. A score from 0-100\n"
                            "2. Detailed feedback on what was good and what could improve\n"
                            "3. An improved version of the answer\n"
                            "4. STAR method analysis (Situation, Task, Action, Result)\n"
                            "5. Specific strengths and weaknesses\n\n"
                            "Be constructive and specific in your feedback."
                        ),
                    },
                    {
                        "role": "user",
                        "content": (
                            f"Job Description:\n{job_description[:2000]}\n\n"
                            f"Interview Question:\n{question}\n\n"
                            f"Candidate's Answer:\n{answer}"
                        ),
                    },
                ],
            )

            return {
                "score": result.score,
                "feedback": result.feedback,
                "improved_answer": result.improved_answer,
                "star_analysis": result.star_analysis.model_dump(),
                "strengths": result.strengths,
                "weaknesses": result.weaknesses,
            }

        except Exception as e:
            logger.error(f"Answer evaluation failed: {e}")
            raise
