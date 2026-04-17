import hashlib
import json
import logging
from typing import Optional

logger = logging.getLogger(__name__)


class LLMCache:
    def __init__(self, redis_url: str = "redis://redis:6379/3"):
        self._redis_url = redis_url
        self._redis = None
        try:
            import redis as redis_lib

            self._redis = redis_lib.Redis.from_url(
                self._redis_url, decode_responses=True
            )
            # Test connection
            self._redis.ping()
        except Exception as e:
            logger.warning(f"Redis cache unavailable: {e}")
            self._redis = None

    @property
    def redis(self):
        return self._redis

    def make_key(self, template_name: str, model: str, input_text: str) -> str:
        """Create a SHA-256 hash key from template name, model, and input text."""
        raw = f"{template_name}:{model}:{input_text}"
        return f"llm_cache:{hashlib.sha256(raw.encode()).hexdigest()}"

    def get(self, prompt_hash: str) -> Optional[str]:
        """Retrieve a cached response by its hash key."""
        try:
            if self.redis is None:
                return None
            return self.redis.get(prompt_hash)
        except Exception as e:
            logger.warning(f"Cache get failed: {e}")
            return None

    def set(self, prompt_hash: str, response: str, ttl: int = 86400) -> None:
        """Store a response in cache with TTL (default 24 hours)."""
        try:
            if self.redis is None:
                return
            self.redis.setex(prompt_hash, ttl, response)
        except Exception as e:
            logger.warning(f"Cache set failed: {e}")

    def invalidate(self, prompt_hash: str) -> None:
        """Remove a cached entry."""
        try:
            if self.redis is None:
                return
            self.redis.delete(prompt_hash)
        except Exception as e:
            logger.warning(f"Cache invalidate failed: {e}")


# Singleton instance for use across the service
_cache_instance: Optional[LLMCache] = None


def get_cache() -> LLMCache:
    """Get or create the singleton LLMCache instance."""
    global _cache_instance
    if _cache_instance is None:
        _cache_instance = LLMCache()
    return _cache_instance
