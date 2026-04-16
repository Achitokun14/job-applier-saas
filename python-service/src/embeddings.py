import logging
import numpy as np

logger = logging.getLogger(__name__)


class EmbeddingService:
    def __init__(self):
        self._model = None
        self._index = None

    @property
    def model(self):
        if self._model is None:
            from sentence_transformers import SentenceTransformer

            self._model = SentenceTransformer("all-MiniLM-L6-v2")
        return self._model

    def embed(self, text: str) -> list[float]:
        """Generate embedding for text."""
        embedding = self.model.encode(text, normalize_embeddings=True)
        return embedding.tolist()

    def compute_match_score(
        self,
        resume_text: str,
        job_description: str,
        resume_skills: list[str] = None,
        job_skills: list[str] = None,
    ) -> dict:
        """Compute weighted match score between resume and job."""
        # 1. Semantic similarity (40%)
        resume_emb = np.array(self.embed(resume_text))
        job_emb = np.array(self.embed(job_description))
        semantic_score = float(np.dot(resume_emb, job_emb)) * 100

        # 2. Skills overlap - Jaccard similarity (30%)
        if resume_skills and job_skills:
            resume_set = set(s.lower().strip() for s in resume_skills)
            job_set = set(s.lower().strip() for s in job_skills)
            if job_set:
                skills_score = len(resume_set & job_set) / len(job_set) * 100
            else:
                skills_score = 50
        else:
            skills_score = 50

        # 3. Experience level heuristic (15%)
        exp_keywords = {
            "senior": 80,
            "lead": 85,
            "principal": 90,
            "junior": 40,
            "intern": 20,
            "mid": 60,
            "staff": 85,
        }
        exp_score = 60
        for kw, score in exp_keywords.items():
            if kw in resume_text.lower() and kw in job_description.lower():
                exp_score = score
                break

        # 4. Education heuristic (15%)
        edu_keywords = {"phd": 90, "master": 75, "bachelor": 60, "mba": 80}
        edu_score = 50
        for kw, score in edu_keywords.items():
            if kw in resume_text.lower():
                edu_score = max(edu_score, score)

        overall = (
            (semantic_score * 0.4)
            + (skills_score * 0.3)
            + (exp_score * 0.15)
            + (edu_score * 0.15)
        )

        return {
            "overall_score": round(min(max(overall, 0), 100)),
            "breakdown": {
                "semantic_similarity": round(semantic_score),
                "skills_match": round(skills_score),
                "experience_match": round(exp_score),
                "education_match": round(edu_score),
            },
        }
