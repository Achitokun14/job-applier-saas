const API_URL = import.meta.env.PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:8080');

interface RequestOptions {
  method?: string;
  body?: unknown;
  token?: string;
}

async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, token } = options;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

export const api = {
  // Auth
  login: (email: string, password: string) =>
    request<{ token: string; user: { id: number; email: string; name: string } }>('/api/v1/auth/login', {
      method: 'POST',
      body: { email, password }
    }),

  register: (email: string, password: string, name: string) =>
    request<{ token: string; user: { id: number; email: string; name: string } }>('/api/v1/auth/register', {
      method: 'POST',
      body: { email, password, name }
    }),

  // Jobs
  searchJobs: (query: string, token: string) =>
    request<Array<{ id: number; title: string; company: string; location: string; description: string }>>(
      `/api/v1/jobs?q=${encodeURIComponent(query)}`,
      { token }
    ),

  applyJob: (jobId: number, token: string) =>
    request<{ message: string }>(`/api/v1/jobs/${jobId}/apply`, {
      method: 'POST',
      token
    }),

  // Applications
  getApplications: (token: string) =>
    request<Array<{ id: number; job: { title: string; company: string }; status: string; applied_at: string }>>(
      '/api/v1/applications',
      { token }
    ),

  getApplication: (id: number, token: string) =>
    request<{ id: number; job: { title: string; company: string }; status: string }>(
      `/api/v1/applications/${id}`,
      { token }
    ),

  deleteApplication: (id: number, token: string) =>
    request<{ message: string }>(`/api/v1/applications/${id}`, {
      method: 'DELETE',
      token
    }),

  // Profile
  getProfile: (token: string) =>
    request<{ id: number; name: string; email: string; resume?: Record<string, unknown> }>('/api/v1/profile', {
      token
    }),

  updateProfile: (data: Record<string, unknown>, token: string) =>
    request<{ message: string }>('/api/v1/profile', {
      method: 'PUT',
      body: data,
      token
    }),

  // Settings
  getSettings: (token: string) =>
    request<{ llm_provider: string; llm_model: string; job_search_remote: boolean; experience_level: string }>(
      '/api/v1/settings',
      { token }
    ),

  updateSettings: (data: Record<string, unknown>, token: string) =>
    request<{ message: string }>('/api/v1/settings', {
      method: 'PUT',
      body: data,
      token
    }),

  // Resume
  generateResume: (resumeYaml: string, style: string, token: string) =>
    request<{ id: string; pdf_path: string }>('/api/v1/resume/generate', {
      method: 'POST',
      body: { resume_yaml: resumeYaml, style },
      token
    }),

  generateCoverLetter: (resumeText: string, jobDescription: string, token: string) =>
    request<{ id: string; pdf_path: string }>('/api/v1/cover-letter/generate', {
      method: 'POST',
      body: { resume_text: resumeText, job_description: jobDescription },
      token
    })
};
