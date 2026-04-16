# Architecture

## System Overview

Auto Job Applier SaaS is a distributed system composed of multiple microservices, each responsible for a specific domain. The architecture follows a clean separation of concerns with well-defined API boundaries.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Layer                                    │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐ │
│  │   Web Browser       │  │   TUI Client        │  │   API Clients       │ │
│  │   (SvelteKit)       │  │   (Bubble Tea)      │  │   (External)        │ │
│  └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘ │
└─────────────┼────────────────────────┼────────────────────────┼─────────────┘
              │                        │                        │
┌─────────────▼────────────────────────▼────────────────────────▼─────────────┐
│                           Gateway Layer                                      │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                         Caddy Reverse Proxy                          │   │
│  │                    Port 80 (HTTP) / 443 (HTTPS)                      │   │
│  └───────────────────────────────┬──────────────────────────────────────┘   │
└──────────────────────────────────┼──────────────────────────────────────────┘
                                   │
┌──────────────────────────────────▼──────────────────────────────────────────┐
│                           Service Layer                                      │
│                                                                              │
│  ┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────┐   │
│  │    Frontend Service  │  │   Backend Service    │  │ Python Service   │   │
│  │    ────────────────  │  │   ────────────────  │  │ ──────────────   │   │
│  │    SvelteKit + Bun   │  │   Go + Chi          │  │ FastAPI          │   │
│  │    Port 3000         │  │   Port 8080         │  │ Port 8001        │   │
│  │                      │  │                      │  │                  │   │
│  │  - Dashboard UI      │  │  - REST API          │  │ - Resume Gen     │   │
│  │  - Auth Pages        │  │  - JWT Auth          │  │ - Cover Letter   │   │
│  │  - Job Search        │  │  - Job Management    │  │ - Job Parser     │   │
│  │  - Applications      │  │  - User Management   │  │ - LLM Integration│   │
│  │  - Profile/Settings  │  │  - Settings          │  │                  │   │
│  └──────────┬───────────┘  └──────────┬───────────┘  └────────┬─────────┘   │
└─────────────┼─────────────────────────┼───────────────────────┼─────────────┘
              │                         │                       │
┌─────────────▼─────────────────────────▼───────────────────────▼─────────────┐
│                           Data Layer                                         │
│                                                                              │
│  ┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────┐   │
│  │    SQLite (Dev)      │  │   PostgreSQL (Prod)  │  │ File Storage     │   │
│  │    ──────────────    │  │   ─────────────────  │  │ ─────────────    │   │
│  │    - Local DB file   │  │   - Managed DB       │  │ - Resumes (PDF)  │   │
│  │    - Auto-migrate    │  │   - Connection pool   │  │ - Cover Letters  │   │
│  │                      │  │   - Full ACID         │  │ - Output files   │   │
│  └──────────────────────┘  └──────────────────────┘  └──────────────────┘   │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. Frontend (SvelteKit + Bun)

```
frontend/
├── src/
│   ├── routes/
│   │   ├── +layout.svelte      # Main layout with navigation
│   │   ├── +page.svelte        # Landing page
│   │   ├── login/              # Authentication
│   │   ├── register/           # User registration
│   │   ├── dashboard/          # Main dashboard
│   │   ├── jobs/               # Job search
│   │   ├── applications/       # Application tracking
│   │   ├── profile/            # User profile
│   │   └── settings/           # Configuration
│   └── lib/
│       ├── api/                # API client
│       └── stores/             # Svelte stores (auth)
├── static/                     # Static assets
└── package.json
```

**Key Patterns:**
- SvelteKit file-based routing
- Svelte stores for state management
- JWT token stored in localStorage
- API calls via fetch with auth headers

### 2. Backend (Go + Chi)

```
backend/
├── cmd/
│   └── server/
│       └── main.go             # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── database/               # Database connection
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data models
│   └── middleware/             # HTTP middleware
└── go.mod
```

**Key Patterns:**
- Chi router with middleware chain
- GORM for database operations
- JWT authentication via middleware
- Repository pattern for data access

### 3. Python Service (FastAPI)

```
python-service/
├── main.py                     # FastAPI app
├── src/
│   ├── resume_generator.py     # Resume generation
│   ├── cover_letter_generator.py # Cover letter generation
│   └── job_parser.py           # Job description parsing
├── data_folder/                # Output files
└── requirements.txt
```

**Key Patterns:**
- FastAPI with Pydantic models
- LangChain for LLM integration
- Selenium for PDF generation
- Fallback templates when LLM unavailable

### 4. TUI (Bubble Tea)

```
tui/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── api/                    # API client
│   ├── models/                 # Data models
│   └── ui/                     # TUI components
└── go.mod
```

**Key Patterns:**
- Bubble Tea Elm architecture
- State machine for navigation
- HTTP client for API calls

## Data Flow

### Authentication Flow

```
User                Frontend              Backend              Database
 │                      │                    │                    │
 │  Login Request       │                    │                    │
 │─────────────────────>│                    │                    │
 │                      │  POST /auth/login  │                    │
 │                      │───────────────────>│                    │
 │                      │                    │  Query user        │
 │                      │                    │───────────────────>│
 │                      │                    │  User data         │
 │                      │                    │<───────────────────│
 │                      │                    │  Verify password   │
 │                      │                    │  Generate JWT      │
 │                      │  JWT + User        │                    │
 │                      │<───────────────────│                    │
 │  Store JWT           │                    │                    │
 │  Redirect            │                    │                    │
 │<─────────────────────│                    │                    │
```

### Job Search Flow

```
User                Frontend              Backend              Database
 │                      │                    │                    │
 │  Search Query        │                    │                    │
 │─────────────────────>│                    │                    │
 │                      │  GET /jobs?q=...   │                    │
 │                      │───────────────────>│                    │
 │                      │                    │  Search jobs       │
 │                      │                    │───────────────────>│
 │                      │                    │  Job results       │
 │                      │                    │<───────────────────│
 │                      │  Job list          │                    │
 │                      │<───────────────────│                    │
 │  Display results     │                    │                    │
 │<─────────────────────│                    │                    │
```

### Resume Generation Flow

```
User                Frontend              Backend         Python Service
 │                      │                    │                    │
 │  Generate Resume     │                    │                    │
 │─────────────────────>│                    │                    │
 │                      │  POST /resume/gen  │                    │
 │                      │───────────────────>│                    │
 │                      │                    │  POST /generate    │
 │                      │                    │───────────────────>│
 │                      │                    │                    │  LLM Processing
 │                      │                    │                    │  PDF Generation
 │                      │                    │  PDF path          │
 │                      │                    │<───────────────────│
 │                      │  PDF download URL  │                    │
 │                      │<───────────────────│                    │
 │  Download PDF        │                    │                    │
 │<─────────────────────│                    │                    │
```

## Database Schema

### Entity Relationship Diagram

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     Users       │       │    Resumes      │       │     Jobs        │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │──┐    │ id (PK)         │       │ id (PK)         │
│ email (UNIQUE)  │  │    │ user_id (FK)    │       │ external_id     │
│ password        │  │    │ personal_info   │       │ title           │
│ name            │  │    │ education       │       │ company         │
│ created_at      │  │    │ experience      │       │ location        │
│ updated_at      │  │    │ skills          │       │ description     │
└─────────────────┘  │    │ projects        │       │ url             │
         │           │    │ pdf_path        │       │ source          │
         │           │    │ created_at      │       │ remote          │
         │           │    │ updated_at      │       │ salary          │
         │           │    └─────────────────┘       │ created_at      │
         │           │              │               └─────────────────┘
         │           │              │                        │
         ▼           │              ▼                        │
┌─────────────────┐  │    ┌─────────────────┐               │
│  Applications   │  │    │   Settings      │               │
├─────────────────┤  │    ├─────────────────┤               │
│ id (PK)         │  │    │ id (PK)         │               │
│ user_id (FK)    │◄─┼────│ user_id (FK)    │               │
│ job_id (FK)     │◄─┼────│ llm_provider    │               │
│ status          │  │    │ llm_model       │               │
│ resume_pdf      │  │    │ llm_api_key     │               │
│ cover_pdf       │  │    │ job_search_*    │               │
│ notes           │  │    │ experience_level│               │
│ applied_at      │  │    │ job_types       │               │
│ created_at      │  │    │ positions       │               │
│ updated_at      │  │    │ locations       │               │
└─────────────────┘  │    │ distance        │               │
                     │    │ created_at      │               │
                     │    │ updated_at      │               │
                     │    └─────────────────┘               │
                     │                                      │
                     └──────────────────────────────────────┘
```

## Service Communication

### Internal API Calls

| From | To | Protocol | Purpose |
|------|-----|----------|---------|
| Frontend | Backend | HTTP/REST | All API operations |
| Backend | Python Service | HTTP/REST | Resume/cover letter generation |
| TUI | Backend | HTTP/REST | All API operations |

### External Integrations

| Service | Provider | Purpose |
|---------|----------|---------|
| LLM | OpenAI, Anthropic, Google | AI-powered content generation |
| Browser | Selenium + Chrome | PDF generation, job scraping |

## Deployment Architecture

### Development

```
Developer Machine
├── Frontend (npm run dev) :5173
├── Backend (go run) :8080
├── Python Service (uvicorn) :8001
├── SQLite ./data/jobapplier.db
└── TUI (go run)
```

### Production (Docker)

```
Docker Host
├── Caddy Container :80/:443
├── Frontend Container :3000
├── Backend Container :8080
├── Python Service Container :8001
├── PostgreSQL Container :5432
└── Volumes
    ├── postgres_data
    ├── caddy_data
    └── data_folder (resumes)
```

## Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Internet                                 │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                    Caddy (TLS Termination)                       │
│                    - HTTPS enforcement                          │
│                    - Security headers                           │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                    Internal Network                              │
│                                                                  │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐       │
│  │   Frontend    │  │   Backend     │  │   Python      │       │
│  │   (no auth)   │  │   (JWT auth)  │  │   (internal)  │       │
│  └───────────────┘  └───────────────┘  └───────────────┘       │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    PostgreSQL                             │   │
│  │                    (no external access)                   │   │
│  └──────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────┘
```

## Scalability Considerations

### Horizontal Scaling

- **Frontend**: Stateless, can run multiple instances behind load balancer
- **Backend**: Stateless, can run multiple instances with shared database
- **Python Service**: Stateless, can run multiple instances
- **Database**: PostgreSQL supports read replicas

### Vertical Scaling

- **Frontend**: Minimal resource requirements
- **Backend**: CPU-bound for request handling
- **Python Service**: Memory-bound for LLM operations
- **Database**: I/O-bound for queries

## Technology Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Frontend Framework | SvelteKit | Fast, modern, excellent DX |
| Frontend Runtime | Bun | Fast JS runtime, native TypeScript |
| Backend Language | Go | Fast, concurrent, simple deployment |
| Backend Router | Chi | Lightweight, idiomatic Go |
| Python Framework | FastAPI | Async, auto-docs, Pydantic validation |
| TUI Framework | Bubble Tea | Mature, well-designed TUI library |
| Database (Dev) | SQLite | Zero-config, file-based |
| Database (Prod) | PostgreSQL | ACID, scalable, feature-rich |
| ORM | GORM | Go idiomatic, multi-database |
| Reverse Proxy | Caddy | Auto-HTTPS, simple config |
| Containerization | Docker | Consistent environments |
