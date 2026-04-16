import type { Handle } from '@sveltejs/kit';

const BACKEND_URL = process.env.BACKEND_URL || 'http://backend:8080';

export const handle: Handle = async ({ event, resolve }) => {
  const path = event.url.pathname;

  // Proxy /api/* and /health to the Go backend
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

      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
      });
    } catch (err) {
      return new Response(JSON.stringify({ error: 'Backend unavailable' }), {
        status: 502,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  }

  return resolve(event);
};
