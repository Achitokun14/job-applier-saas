import type { Handle } from '@sveltejs/kit';

const BACKEND_URL = process.env.BACKEND_URL || 'http://backend:8080';

const SECURITY_HEADERS: Record<string, string> = {
  'X-Frame-Options': 'DENY',
  'X-Content-Type-Options': 'nosniff',
  'Referrer-Policy': 'strict-origin-when-cross-origin',
  'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
  'Content-Security-Policy': "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https://jobs.murgana.online",
};

export const handle: Handle = async ({ event, resolve }) => {
  const path = event.url.pathname;

  // Handle CORS preflight
  if (event.request.method === 'OPTIONS') {
    return new Response(null, {
      status: 204,
      headers: {
        'Access-Control-Allow-Origin': 'https://jobs.murgana.online',
        'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization',
        'Access-Control-Max-Age': '86400',
      },
    });
  }

  // Proxy /api/* and /health and /metrics to the Go backend
  if (path.startsWith('/api/') || path === '/health' || path === '/metrics') {
    const targetUrl = `${BACKEND_URL}${path}${event.url.search}`;

    const headers = new Headers();
    for (const [key, value] of event.request.headers) {
      if (key.toLowerCase() !== 'host') {
        headers.set(key, value);
      }
    }

    try {
      const response = await fetch(targetUrl, {
        method: event.request.method,
        headers,
        body: event.request.method !== 'GET' && event.request.method !== 'HEAD'
          ? await event.request.text()
          : undefined,
      });

      // Build response headers - copy from backend + add security headers
      const responseHeaders = new Headers();
      for (const [key, value] of response.headers) {
        responseHeaders.set(key, value);
      }
      for (const [key, value] of Object.entries(SECURITY_HEADERS)) {
        responseHeaders.set(key, value);
      }

      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: responseHeaders,
      });
    } catch (err) {
      return new Response(JSON.stringify({ error: 'Backend unavailable' }), {
        status: 502,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  }

  // For frontend pages, add security headers
  const response = await resolve(event);
  const newHeaders = new Headers(response.headers);
  for (const [key, value] of Object.entries(SECURITY_HEADERS)) {
    newHeaders.set(key, value);
  }

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: newHeaders,
  });
};
