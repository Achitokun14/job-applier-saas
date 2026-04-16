import logging
import re
from pathlib import Path
from urllib.parse import urlparse
from pydantic import BaseModel, Field
from typing import Optional

logger = logging.getLogger(__name__)

CHROMA_DIR = Path("data_folder/chroma")
CHROMA_DIR.mkdir(parents=True, exist_ok=True)


class CompanyInfo(BaseModel):
    company_info: str = Field(
        default="", description="General overview of the company"
    )
    culture: str = Field(
        default="", description="Company culture and values"
    )
    recent_news: str = Field(
        default="", description="Recent news, developments, or announcements"
    )
    key_facts: list[str] = Field(
        default_factory=list,
        description="Key facts about the company useful for interviews/cover letters",
    )
    industry: str = Field(default="", description="Industry the company operates in")
    size: str = Field(default="", description="Company size estimate")
    mission: str = Field(default="", description="Company mission statement or purpose")


class CompanyResearcher:
    def __init__(self):
        self._chroma_client = None

    @property
    def chroma_client(self):
        if self._chroma_client is None:
            import chromadb

            self._chroma_client = chromadb.PersistentClient(
                path=str(CHROMA_DIR)
            )
        return self._chroma_client

    def _fetch_page(self, url: str) -> str:
        """Fetch a web page and extract text content."""
        try:
            import httpx

            headers = {
                "User-Agent": (
                    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
                    "AppleWebKit/537.36 (KHTML, like Gecko) "
                    "Chrome/120.0.0.0 Safari/537.36"
                )
            }
            response = httpx.get(url, headers=headers, timeout=15, follow_redirects=True)
            response.raise_for_status()

            # Extract text using selectolax (already a dependency)
            from selectolax.parser import HTMLParser

            tree = HTMLParser(response.text)

            # Remove script and style elements
            for tag in tree.css("script, style, nav, footer, header"):
                tag.decompose()

            text = tree.body.text(separator="\n") if tree.body else ""
            # Clean up whitespace
            text = re.sub(r"\n{3,}", "\n\n", text)
            text = re.sub(r" {2,}", " ", text)
            return text.strip()

        except Exception as e:
            logger.warning(f"Failed to fetch {url}: {e}")
            return ""

    def _chunk_text(self, text: str, chunk_size: int = 500, overlap: int = 50) -> list[str]:
        """Split text into overlapping chunks."""
        if not text:
            return []
        words = text.split()
        chunks = []
        for i in range(0, len(words), chunk_size - overlap):
            chunk = " ".join(words[i : i + chunk_size])
            if chunk.strip():
                chunks.append(chunk)
        return chunks

    def _get_collection_name(self, company_name: str) -> str:
        """Generate a valid ChromaDB collection name from company name."""
        # ChromaDB collection names must be 3-63 chars, alphanumeric with underscores/hyphens
        name = re.sub(r"[^a-zA-Z0-9_-]", "_", company_name.lower().strip())
        name = re.sub(r"_+", "_", name).strip("_")
        if len(name) < 3:
            name = name + "_co"
        return name[:63]

    def _store_in_chroma(self, company_name: str, texts: list[str]) -> None:
        """Store text chunks in ChromaDB collection."""
        if not texts:
            return

        collection_name = self._get_collection_name(company_name)
        try:
            collection = self.chroma_client.get_or_create_collection(
                name=collection_name
            )

            # Add documents with unique IDs
            ids = [f"{collection_name}_{i}" for i in range(len(texts))]
            collection.upsert(documents=texts, ids=ids)
            logger.info(
                f"Stored {len(texts)} chunks for company '{company_name}' "
                f"in collection '{collection_name}'"
            )
        except Exception as e:
            logger.warning(f"ChromaDB storage failed: {e}")

    def _retrieve_context(self, company_name: str, query: str, n_results: int = 5) -> str:
        """Retrieve relevant context from ChromaDB."""
        collection_name = self._get_collection_name(company_name)
        try:
            collection = self.chroma_client.get_collection(name=collection_name)
            results = collection.query(query_texts=[query], n_results=n_results)
            if results and results["documents"]:
                return "\n\n".join(results["documents"][0])
        except Exception as e:
            logger.warning(f"ChromaDB retrieval failed: {e}")
        return ""

    def research(
        self,
        company_name: str,
        company_url: str = None,
        llm_provider: str = None,
        llm_api_key: str = None,
        llm_model: str = None,
    ) -> dict:
        """Research a company using web scraping + RAG + LLM analysis."""
        fetched_texts = []

        # Fetch company website pages
        if company_url:
            # Normalize URL
            if not company_url.startswith("http"):
                company_url = f"https://{company_url}"

            parsed = urlparse(company_url)
            base_url = f"{parsed.scheme}://{parsed.netloc}"

            # Try fetching key pages
            pages_to_fetch = [
                company_url,
                f"{base_url}/about",
                f"{base_url}/about-us",
                f"{base_url}/careers",
                f"{base_url}/culture",
                f"{base_url}/values",
            ]

            for page_url in pages_to_fetch:
                text = self._fetch_page(page_url)
                if text and len(text) > 100:
                    fetched_texts.append(text[:3000])  # Limit per page

        # Chunk and store in ChromaDB
        all_chunks = []
        for text in fetched_texts:
            all_chunks.extend(self._chunk_text(text))

        if all_chunks:
            self._store_in_chroma(company_name, all_chunks)

        # Retrieve relevant context
        rag_context = self._retrieve_context(
            company_name,
            f"company culture values mission {company_name}",
        )

        # Use LLM to synthesize company research
        try:
            from src.llm_client import get_client

            client, model = get_client(
                provider=llm_provider, api_key=llm_api_key, model=llm_model
            )

            web_context = "\n\n".join(fetched_texts[:3]) if fetched_texts else ""
            combined_context = f"{web_context}\n\n{rag_context}".strip()

            if combined_context:
                prompt = (
                    f"Based on the following information about {company_name}, "
                    f"provide a comprehensive company research summary:\n\n"
                    f"{combined_context[:6000]}\n\n"
                    "Synthesize the information into a structured company profile."
                )
            else:
                prompt = (
                    f"Provide a comprehensive research summary about the company "
                    f"'{company_name}'. Include general overview, culture, recent "
                    f"developments, key facts, industry, size, and mission. "
                    f"Use your training knowledge."
                )

            result = client.chat.completions.create(
                model=model,
                response_model=CompanyInfo,
                max_retries=2,
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are a company research analyst. Provide accurate, "
                            "detailed information about companies to help job seekers "
                            "prepare for interviews and write cover letters. "
                            "Be factual and specific."
                        ),
                    },
                    {"role": "user", "content": prompt},
                ],
            )

            return result.model_dump()

        except Exception as e:
            logger.error(f"Company research failed: {e}")
            raise
