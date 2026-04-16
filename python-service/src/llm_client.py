import json
import logging
import os
from typing import Optional

import instructor
import litellm

logger = logging.getLogger(__name__)

# Map provider names to LiteLLM model prefixes
_PROVIDER_PREFIXES = {
    "openai": "",
    "anthropic": "anthropic/",
    "google": "gemini/",
    "groq": "groq/",
    "mistral": "mistral/",
    "ollama": "ollama/",
}

# Default models per provider
_DEFAULT_MODELS = {
    "openai": "gpt-4o-mini",
    "anthropic": "claude-sonnet-4-20250514",
    "google": "gemini-1.5-flash",
    "groq": "llama-3.1-8b-instant",
    "mistral": "mistral-small-latest",
    "ollama": "llama3",
}


def get_client(
    provider: Optional[str] = None,
    api_key: Optional[str] = None,
    model: Optional[str] = None,
):
    """
    Return a tuple of (instructor_client, full_model_name) ready for
    ``client.chat.completions.create(model=..., response_model=...)``.

    Falls back to the environment variables ``LLM_PROVIDER``,
    ``LLM_API_KEY``, and ``LLM_MODEL`` when arguments are not supplied.

    Supported providers: openai, anthropic, google, groq, mistral, ollama.
    """
    provider = (provider or os.getenv("LLM_PROVIDER", "openai")).lower().strip()
    api_key = api_key or os.getenv("LLM_API_KEY", "")
    model = model or os.getenv("LLM_MODEL", "")

    if not model:
        model = _DEFAULT_MODELS.get(provider, "gpt-4o-mini")

    prefix = _PROVIDER_PREFIXES.get(provider, "")
    # Only add the prefix if the model string doesn't already include it
    if prefix and not model.startswith(prefix):
        full_model = f"{prefix}{model}"
    else:
        full_model = model

    # Set the API key in the environment so LiteLLM picks it up.
    # Each provider has its own env-var name, but we also accept a
    # generic LLM_API_KEY for convenience.
    if api_key:
        _env_key_map = {
            "openai": "OPENAI_API_KEY",
            "anthropic": "ANTHROPIC_API_KEY",
            "google": "GEMINI_API_KEY",
            "groq": "GROQ_API_KEY",
            "mistral": "MISTRAL_API_KEY",
        }
        env_var = _env_key_map.get(provider)
        if env_var:
            os.environ[env_var] = api_key

    client = instructor.from_litellm(litellm.completion)

    return client, full_model


def cached_llm_call(
    template_name: str,
    messages: list[dict],
    response_model,
    provider: Optional[str] = None,
    api_key: Optional[str] = None,
    model: Optional[str] = None,
    max_retries: int = 2,
    cache_ttl: int = 86400,
):
    """
    Call the LLM with caching support. Checks Redis cache before calling
    the LLM and stores the response after.

    Returns the parsed Pydantic model instance.
    """
    from src.llm_cache import get_cache

    client, full_model = get_client(provider=provider, api_key=api_key, model=model)
    cache = get_cache()

    # Build cache key from template, model, and message content
    input_text = json.dumps(messages, sort_keys=True)
    cache_key = cache.make_key(template_name, full_model, input_text)

    # Check cache
    cached = cache.get(cache_key)
    if cached is not None:
        try:
            data = json.loads(cached)
            return response_model.model_validate(data)
        except Exception as e:
            logger.warning(f"Cache deserialization failed: {e}")

    # Call LLM
    result = client.chat.completions.create(
        model=full_model,
        response_model=response_model,
        max_retries=max_retries,
        messages=messages,
    )

    # Store in cache
    try:
        cache.set(cache_key, result.model_dump_json(), ttl=cache_ttl)
    except Exception as e:
        logger.warning(f"Failed to cache LLM response: {e}")

    return result
