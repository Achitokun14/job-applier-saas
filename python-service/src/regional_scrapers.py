import hashlib
import logging
from typing import Optional

import httpx
from selectolax.parser import HTMLParser

logger = logging.getLogger(__name__)

USER_AGENT = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/124.0.0.0 Safari/537.36"
)


def _generate_id(source: str, title: str, company: str) -> str:
    raw = f"{source}-{title}-{company}"
    return hashlib.sha256(raw.encode()).hexdigest()[:16]


def _text(node) -> str:
    """Safely extract text from a selectolax node."""
    if node is None:
        return ""
    return node.text(strip=True)


class RegionalScraper:
    """Supplemental scrapers for regional job boards not covered by JobSpy."""

    def __init__(self):
        self._headers = {"User-Agent": USER_AGENT}

    async def _fetch(self, url: str) -> Optional[HTMLParser]:
        try:
            async with httpx.AsyncClient(
                headers=self._headers, timeout=15.0, follow_redirects=True
            ) as client:
                resp = await client.get(url)
                resp.raise_for_status()
                return HTMLParser(resp.text)
        except Exception as e:
            logger.error(f"Failed to fetch {url}: {e}")
            return None

    async def scrape_rekrute(
        self, search_term: str, location: Optional[str] = None
    ) -> list[dict]:
        """Scrape rekrute.com (Morocco / North Africa)."""
        try:
            query = search_term.replace(" ", "+")
            url = f"https://www.rekrute.com/offres.html?s=3&p=1&o=1&keyword={query}"
            if location:
                url += f"&l={location.replace(' ', '+')}"

            tree = await self._fetch(url)
            if tree is None:
                return []

            jobs = []
            for card in tree.css("li.post-id"):
                title_el = card.css_first("h2 a, .titreJob a, a.titreJob")
                company_el = card.css_first(
                    ".info span, .entreprise, .holder .info a"
                )
                location_el = card.css_first(
                    ".info span:nth-child(2), .localisation"
                )
                link_el = card.css_first("h2 a, a.titreJob")

                title = _text(title_el)
                company = _text(company_el)
                loc = _text(location_el) or (location or "Morocco")
                href = link_el.attributes.get("href", "") if link_el else ""
                job_url = href if href.startswith("http") else f"https://www.rekrute.com{href}"

                if not title:
                    continue

                jobs.append(
                    {
                        "external_id": _generate_id("rekrute", title, company),
                        "source": "rekrute",
                        "title": title,
                        "company": company,
                        "location": loc,
                        "description": "",
                        "url": job_url,
                        "remote": False,
                        "salary": "",
                        "posted_at": "",
                    }
                )
            return jobs
        except Exception as e:
            logger.error(f"Rekrute scrape failed: {e}")
            return []

    async def scrape_careerlink(
        self, search_term: str, location: Optional[str] = None
    ) -> list[dict]:
        """Scrape careerlink.ma (Morocco)."""
        try:
            query = search_term.replace(" ", "+")
            url = f"https://www.careerlink.ma/emploi?q={query}"
            if location:
                url += f"&l={location.replace(' ', '+')}"

            tree = await self._fetch(url)
            if tree is None:
                return []

            jobs = []
            for card in tree.css(
                "article, .job-card, .job-listing, .job-item, .listing-item"
            ):
                title_el = card.css_first("h2 a, h3 a, .job-title a, a.title")
                company_el = card.css_first(
                    ".company, .company-name, .employer"
                )
                location_el = card.css_first(".location, .job-location, .city")
                link_el = card.css_first("h2 a, h3 a, .job-title a, a.title")

                title = _text(title_el)
                company = _text(company_el)
                loc = _text(location_el) or (location or "Morocco")
                href = link_el.attributes.get("href", "") if link_el else ""
                job_url = (
                    href
                    if href.startswith("http")
                    else f"https://www.careerlink.ma{href}"
                )

                if not title:
                    continue

                jobs.append(
                    {
                        "external_id": _generate_id("careerlink", title, company),
                        "source": "careerlink",
                        "title": title,
                        "company": company,
                        "location": loc,
                        "description": "",
                        "url": job_url,
                        "remote": False,
                        "salary": "",
                        "posted_at": "",
                    }
                )
            return jobs
        except Exception as e:
            logger.error(f"Careerlink scrape failed: {e}")
            return []

    async def scrape_emploima(
        self, search_term: str, location: Optional[str] = None
    ) -> list[dict]:
        """Scrape emploi-ma.com (Morocco)."""
        try:
            query = search_term.replace(" ", "+")
            url = f"https://www.emploi-ma.com/recherche?q={query}"
            if location:
                url += f"&l={location.replace(' ', '+')}"

            tree = await self._fetch(url)
            if tree is None:
                return []

            jobs = []
            for card in tree.css(
                "article, .job-card, .job-listing, .result-item, .annonce"
            ):
                title_el = card.css_first(
                    "h2 a, h3 a, .job-title a, a.title, .annonce-title a"
                )
                company_el = card.css_first(
                    ".company, .company-name, .employer, .recruteur"
                )
                location_el = card.css_first(
                    ".location, .job-location, .ville, .city"
                )
                link_el = card.css_first(
                    "h2 a, h3 a, .job-title a, a.title, .annonce-title a"
                )

                title = _text(title_el)
                company = _text(company_el)
                loc = _text(location_el) or (location or "Morocco")
                href = link_el.attributes.get("href", "") if link_el else ""
                job_url = (
                    href
                    if href.startswith("http")
                    else f"https://www.emploi-ma.com{href}"
                )

                if not title:
                    continue

                jobs.append(
                    {
                        "external_id": _generate_id("emploima", title, company),
                        "source": "emploima",
                        "title": title,
                        "company": company,
                        "location": loc,
                        "description": "",
                        "url": job_url,
                        "remote": False,
                        "salary": "",
                        "posted_at": "",
                    }
                )
            return jobs
        except Exception as e:
            logger.error(f"EmploiMa scrape failed: {e}")
            return []

    async def scrape_jobberman(
        self, search_term: str, location: Optional[str] = None
    ) -> list[dict]:
        """Scrape jobberman.com (West Africa -- Nigeria, Ghana, etc.)."""
        try:
            query = search_term.replace(" ", "+")
            url = f"https://www.jobberman.com/jobs?q={query}"
            if location:
                url += f"&l={location.replace(' ', '+')}"

            tree = await self._fetch(url)
            if tree is None:
                return []

            jobs = []
            for card in tree.css(
                "article, .job-card, .job-listing, .search-result, .job_listing"
            ):
                title_el = card.css_first(
                    "h2 a, h3 a, .job-title a, a.title, .job__title a"
                )
                company_el = card.css_first(
                    ".company, .company-name, .employer, .job__company"
                )
                location_el = card.css_first(
                    ".location, .job-location, .city, .job__location"
                )
                link_el = card.css_first(
                    "h2 a, h3 a, .job-title a, a.title, .job__title a"
                )

                title = _text(title_el)
                company = _text(company_el)
                loc = _text(location_el) or (location or "")
                href = link_el.attributes.get("href", "") if link_el else ""
                job_url = (
                    href
                    if href.startswith("http")
                    else f"https://www.jobberman.com{href}"
                )

                if not title:
                    continue

                jobs.append(
                    {
                        "external_id": _generate_id("jobberman", title, company),
                        "source": "jobberman",
                        "title": title,
                        "company": company,
                        "location": loc,
                        "description": "",
                        "url": job_url,
                        "remote": False,
                        "salary": "",
                        "posted_at": "",
                    }
                )
            return jobs
        except Exception as e:
            logger.error(f"Jobberman scrape failed: {e}")
            return []

    async def scrape_all(
        self, search_term: str, location: Optional[str] = None
    ) -> list[dict]:
        """Run all regional scrapers and return combined results."""
        import asyncio

        results = await asyncio.gather(
            self.scrape_rekrute(search_term, location),
            self.scrape_careerlink(search_term, location),
            self.scrape_emploima(search_term, location),
            self.scrape_jobberman(search_term, location),
            return_exceptions=True,
        )

        combined = []
        for result in results:
            if isinstance(result, Exception):
                logger.error(f"Regional scraper failed: {result}")
                continue
            combined.extend(result)
        return combined
