# API Reference

Complete REST API documentation for Auto Job Applier SaaS.

## Base URLs

| Environment | URL |
|-------------|-----|
| Development | `http://localhost:8080` |
| Production | `https://your-domain.com` |

## Authentication

All authenticated endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### Register User

```http
POST /api/v1/auth/register
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `409` - Email already registered

---

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid credentials

---

### Jobs

#### Search Jobs

```http
GET /api/v1/jobs?q={query}&source={source}&page={page}
Authorization: Bearer <token>
```

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | No | Search query |
| `source` | string | No | Job source filter (linkedin, indeed, etc.) |
| `page` | integer | No | Page number (default: 1) |

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "external_id": "job-123",
    "title": "Software Engineer",
    "company": "Tech Corp",
    "location": "San Francisco, CA",
    "description": "We are looking for...",
    "url": "https://linkedin.com/jobs/123",
    "source": "linkedin",
    "remote": true,
    "salary": "$120k-$150k"
  }
]
```

---

#### Apply to Job

```http
POST /api/v1/jobs/{id}/apply
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Job ID |

**Response (201 Created):**
```json
{
  "id": 1,
  "user_id": 1,
  "job_id": 1,
  "status": "applied",
  "applied_at": "2026-03-27T10:30:00Z"
}
```

**Errors:**
- `404` - Job not found
- `409` - Already applied to this job

---

### Applications

#### List Applications

```http
GET /api/v1/applications
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "job_id": 1,
    "job": {
      "id": 1,
      "title": "Software Engineer",
      "company": "Tech Corp",
      "location": "San Francisco, CA"
    },
    "status": "applied",
    "resume_pdf": "/output/resume_1.pdf",
    "cover_pdf": "/output/cover_1.pdf",
    "applied_at": "2026-03-27T10:30:00Z"
  }
]
```

---

#### Get Application

```http
GET /api/v1/applications/{id}
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Application ID |

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "job_id": 1,
  "job": {
    "id": 1,
    "title": "Software Engineer",
    "company": "Tech Corp",
    "description": "We are looking for..."
  },
  "status": "applied",
  "resume_pdf": "/output/resume_1.pdf",
  "cover_pdf": "/output/cover_1.pdf",
  "notes": "",
  "applied_at": "2026-03-27T10:30:00Z"
}
```

**Errors:**
- `404` - Application not found

---

#### Delete Application

```http
DELETE /api/v1/applications/{id}
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Application ID |

**Response (200 OK):**
```json
{
  "message": "Application deleted"
}
```

**Errors:**
- `404` - Application not found

---

### Profile

#### Get Profile

```http
GET /api/v1/profile
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "resume": {
    "user_id": 1,
    "personal_info": "{\"name\":\"John Doe\",\"email\":\"john@example.com\"}",
    "education": "[{\"degree\":\"BS Computer Science\",\"institution\":\"MIT\"}]",
    "experience": "[{\"position\":\"Software Engineer\",\"company\":\"Tech Corp\"}]",
    "skills": "[\"Go\",\"Python\",\"JavaScript\"]",
    "projects": "[{\"name\":\"Auto Job Applier\",\"description\":\"...\"}]",
    "pdf_path": "/output/resume.pdf"
  }
}
```

---

#### Update Profile

```http
PUT /api/v1/profile
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "John Doe",
  "personal_info": "{\"name\":\"John Doe\",\"email\":\"john@example.com\"}",
  "education": "[{\"degree\":\"BS Computer Science\",\"institution\":\"MIT\"}]",
  "experience": "[{\"position\":\"Software Engineer\",\"company\":\"Tech Corp\"}]",
  "skills": "[\"Go\",\"Python\",\"JavaScript\"]",
  "projects": "[{\"name\":\"Auto Job Applier\"}]",
  "achievements": "[\"Led team of 5 engineers\"]",
  "certifications": "[\"AWS Solutions Architect\"]",
  "languages": "[\"English\",\"Spanish\"]"
}
```

**Response (200 OK):**
```json
{
  "message": "Profile updated"
}
```

---

### Settings

#### Get Settings

```http
GET /api/v1/settings
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "llm_provider": "openai",
  "llm_model": "gpt-4o-mini",
  "job_search_remote": true,
  "job_search_hybrid": true,
  "job_search_onsite": false,
  "experience_level": "mid_senior",
  "job_types": "full_time",
  "positions": "Software Engineer, Backend Developer",
  "locations": "San Francisco, Remote",
  "distance": 50
}
```

---

#### Update Settings

```http
PUT /api/v1/settings
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "llm_provider": "openai",
  "llm_model": "gpt-4o-mini",
  "llm_api_key": "sk-...",
  "job_search_remote": true,
  "job_search_hybrid": true,
  "job_search_onsite": false,
  "experience_level": "mid_senior",
  "job_types": "full_time",
  "positions": "Software Engineer, Backend Developer",
  "locations": "San Francisco, Remote",
  "distance": 50
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "llm_provider": "openai",
  "llm_model": "gpt-4o-mini",
  "job_search_remote": true
}
```

---

### Resume Generation

#### Generate Resume

```http
POST /api/v1/resume/generate
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "resume_yaml": "personal_information:\n  name: John Doe\n  email: john@example.com",
  "style": "modern",
  "job_url": "https://linkedin.com/jobs/123"
}
```

**Response (200 OK):**
```json
{
  "message": "Resume generation endpoint - connects to Python service",
  "status": "pending"
}
```

---

#### Generate Cover Letter

```http
POST /api/v1/cover-letter/generate
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "resume_text": "John Doe - Software Engineer with 5 years experience...",
  "job_description": "We are looking for a senior software engineer..."
}
```

**Response (200 OK):**
```json
{
  "message": "Cover letter generation endpoint - connects to Python service",
  "status": "pending"
}
```

---

## Python Service API

The Python service runs on port 8001 and provides resume generation capabilities.

### Health Check

```http
GET http://localhost:8001/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "resume-generator"
}
```

---

### Generate Resume

```http
POST http://localhost:8001/generate-resume
Content-Type: application/json
```

**Request Body:**
```json
{
  "resume_yaml": "personal_information:\n  name: John Doe",
  "style": "modern",
  "job_url": "https://linkedin.com/jobs/123",
  "job_description": "Optional job description text"
}
```

**Response:**
```json
{
  "id": "abc12345",
  "pdf_path": "/output/resume_abc12345.pdf",
  "html_content": "<html>...</html>",
  "metadata": {
    "style": "modern",
    "tailored": true,
    "word_count": 500
  }
}
```

---

### Generate Cover Letter

```http
POST http://localhost:8001/generate-cover-letter
Content-Type: application/json
```

**Request Body:**
```json
{
  "resume_text": "John Doe - Software Engineer...",
  "job_description": "We are looking for...",
  "company_name": "Tech Corp",
  "job_title": "Senior Engineer"
}
```

**Response:**
```json
{
  "id": "def67890",
  "pdf_path": "/output/cover_letter_def67890.pdf",
  "html_content": "<html>...</html>",
  "metadata": {
    "company": "Tech Corp",
    "job_title": "Senior Engineer",
    "word_count": 350
  }
}
```

---

### Parse Job

```http
POST http://localhost:8001/parse-job
Content-Type: application/json
```

**Request Body:**
```json
{
  "url": "https://linkedin.com/jobs/123"
}
```

**Response:**
```json
{
  "title": "Software Engineer",
  "company": "Tech Corp",
  "location": "San Francisco, CA",
  "description": "We are looking for...",
  "requirements": ["5+ years experience", "Go proficiency"],
  "responsibilities": ["Design systems", "Lead team"],
  "salary": "$120k-$150k",
  "remote": true
}
```

---

### List Styles

```http
GET http://localhost:8001/styles
```

**Response:**
```json
{
  "styles": [
    {"id": "modern", "name": "Modern", "description": "Clean, contemporary design"},
    {"id": "classic", "name": "Classic", "description": "Traditional, professional"},
    {"id": "minimal", "name": "Minimal", "description": "Simple, elegant"},
    {"id": "creative", "name": "Creative", "description": "Bold design"},
    {"id": "professional", "name": "Professional", "description": "Corporate-focused"}
  ]
}
```

---

### List Templates

```http
GET http://localhost:8001/templates
```

**Response:**
```json
{
  "templates": [
    {"id": "standard", "name": "Standard Resume", "sections": ["experience", "education", "skills"]},
    {"id": "executive", "name": "Executive Resume", "sections": ["summary", "achievements", "experience"]},
    {"id": "technical", "name": "Technical Resume", "sections": ["skills", "projects", "experience"]}
  ]
}
```

---

## Error Responses

All endpoints may return these error formats:

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 409 Conflict
```json
{
  "error": "Resource already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Rate Limiting

- Authentication endpoints: 10 requests per minute
- Other endpoints: 100 requests per minute

## CORS

Allowed origins:
- `http://localhost:5173` (development)
- `http://localhost:3000` (production)
- Configured via `CORS_ALLOWED_ORIGINS` environment variable
