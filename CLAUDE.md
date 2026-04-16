# CLAUDE.md - Job Applier SaaS

## Project Overview
Automated job application SaaS platform. AI-powered resume/cover letter generation, multi-platform job scraping, and automated applications with subscription billing.

## Architecture

| Service | Tech | Port | Purpose |
|---------|------|------|---------|
| **Frontend** | SvelteKit + Bun | 3000 | Web dashboard |
| **Backend API** | Go + Chi | 8080 | REST API, auth, payments, task orchestration |
| **Python Service** | FastAPI | 8001 | AI resume/cover letter generation (LLM) |
| **Celery Worker** | Python + Celery | -- | Async task processing (resume gen, scraping) |
| **Flower** | Celery Flower | 5555 | Celery task monitoring dashboard |
| **Scraper** | Go | 8081 | Job scraping service |
| **PostgreSQL** | Postgres 15 | 5432 | Primary database |
| **Redis** | Redis 7 | 6379 | Cache, Asynq task queue, Celery broker, rate limiting |
| **Asynq Monitor** | asynqmon | 8082 | Asynq task queue dashboard |
| **Caddy** | Caddy 2 | 80/443 | Reverse proxy, HTTPS termination |
| **TUI** | Go + Bubble Tea | -- | Terminal interface for local testing |

## Directory Structure
```
job-applier-saas/
├── backend/
│   ├── cmd/
│   │   ├── server/         # API server entrypoint
│   │   └── migrate/        # Migration CLI
│   ├── internal/
│   │   ├── auth/           # JWT, password hashing
│   │   ├── cache/          # Redis cache layer
│   │   ├── config/         # App configuration
│   │   ├── crypto/         # AES-256-GCM encryption
│   │   ├── database/       # DB connection, GORM setup
│   │   ├── errors/         # Structured error handling
│   │   ├── handlers/       # HTTP handlers (auth, payment, health)
│   │   ├── logger/         # Structured logging
│   │   ├── metrics/        # Prometheus metrics
│   │   ├── middleware/      # Rate limiting, RBAC, security, subscription
│   │   ├── models/         # GORM models
│   │   ├── repository/     # Data access (users, jobs, applications, settings)
│   │   ├── services/       # Business logic (python client, usage tracking)
│   │   └── tasks/          # Asynq async tasks (auto-apply, resume, cover letter, scrape)
│   └── migrations/         # SQL migration files
├── frontend/               # SvelteKit + Bun
├── python-service/         # FastAPI + Celery workers
├── scraper/                # Go job scraper
├── tui/                    # Bubble Tea TUI
├── scripts/
│   ├── backup.sh           # Database backup (pg_dump)
│   ├── restore.sh          # Database restore
│   └── ...                 # Deploy, setup, start/stop scripts
├── docker-compose.yml      # All services
├── docker-compose.production.yml
├── Caddyfile
├── .github/workflows/ci.yml  # CI pipeline
└── .env.example
```

## Development Commands

### Backend (Go)
```bash
cd backend
go run cmd/server/main.go                      # Run API server
go run cmd/migrate/main.go up                   # Run migrations
go run cmd/migrate/main.go down                 # Rollback migration
go test -v -race ./...                          # Run tests with race detector
go build -o bin/server cmd/server/main.go       # Build binary
```

### Frontend (SvelteKit + Bun)
```bash
cd frontend
bun install          # Install dependencies
bun run dev          # Dev server (hot reload)
bun run build        # Production build
bun run check        # Type checking
bun run lint         # Linting
```

### Python Service
```bash
cd python-service
pip install -r requirements.txt
uvicorn main:app --reload --port 8001           # Dev server
celery -A src.celery_app worker --loglevel=info # Celery worker
celery -A src.celery_app flower --port=5555     # Flower monitoring
python -m pytest tests/ -v                      # Run tests
```

### TUI
```bash
cd tui
go run cmd/main.go                              # Run TUI
```

### Docker
```bash
docker compose up -d                  # Start all services
docker compose down                   # Stop all services
docker compose build                  # Rebuild all images
docker compose build backend          # Rebuild single service
docker compose logs -f backend        # Follow backend logs
docker compose logs -f celery-worker  # Follow Celery logs
docker compose ps                     # Service status
docker compose exec postgres psql -U jobapplier  # DB shell
```

### Database Backup/Restore
```bash
./scripts/backup.sh [output_dir]               # Backup (default: ./backups/)
./scripts/restore.sh backups/jobapplier_*.sql.gz  # Restore from backup
```

### Monitoring Dashboards
- Flower (Celery tasks): http://localhost:5555
- Asynq Monitor (Go tasks): http://localhost:8082

## API Endpoints

### Public (no auth)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Register user |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh JWT token |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password |
| POST | `/api/v1/payments/webhook` | Stripe webhook |

### Protected (JWT required)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/logout` | Logout (invalidate token) |
| GET | `/api/v1/jobs` | Search jobs |
| POST | `/api/v1/jobs/{id}/apply` | Apply to job |
| POST | `/api/v1/jobs/{id}/auto-apply` | Auto-apply with AI |
| POST | `/api/v1/jobs/ingest` | Ingest job listings |
| GET | `/api/v1/applications` | List applications |
| GET | `/api/v1/applications/{id}` | Get application details |
| DELETE | `/api/v1/applications/{id}` | Delete application |
| POST | `/api/v1/resume/generate` | Generate tailored resume |
| POST | `/api/v1/cover-letter/generate` | Generate cover letter |
| GET | `/api/v1/tasks/{id}` | Get async task status |
| GET | `/api/v1/profile` | Get user profile |
| PUT | `/api/v1/profile` | Update user profile |
| GET | `/api/v1/settings` | Get user settings |
| PUT | `/api/v1/settings` | Update user settings |
| POST | `/api/v1/scrape/trigger` | Trigger job scraping |
| POST | `/api/v1/payments/checkout` | Create Stripe checkout |
| GET | `/api/v1/payments/subscription` | Get subscription info |
| GET | `/api/v1/payments/portal` | Stripe billing portal |
| POST | `/api/v1/payments/cancel` | Cancel subscription |

## Environment Variables

See `.env.example` for all variables. Key ones:

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | Min 32 chars, change in production |
| `REDIS_URL` | Yes | Redis host:port |
| `LLM_API_KEY` | For AI features | OpenAI/Anthropic API key |
| `LLM_MODEL` | No | Default: `gpt-4o-mini` |
| `STRIPE_SECRET_KEY` | For payments | Stripe secret key |
| `STRIPE_WEBHOOK_SECRET` | For payments | Stripe webhook signing secret |
| `ENCRYPTION_KEY` | Yes | 32+ chars for AES-256-GCM |
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated origins |

## Testing
```bash
# Backend (with race detector and coverage)
cd backend && go test -v -race -coverprofile=coverage.out ./...

# Frontend type checking
cd frontend && bun run check

# Python service
cd python-service && python -m pytest tests/ -v

# Full CI locally
docker compose build
```

## Key Dependencies

**Backend (Go):** Chi (router), GORM (ORM), go-redis, Asynq (task queue), golang-jwt, stripe-go, golang-migrate, gosec (security)

**Frontend:** SvelteKit, Bun

**Python Service:** FastAPI, Celery, Redis (broker)

**Infrastructure:** PostgreSQL 15, Redis 7, Caddy 2, Docker Compose

## Worktree Directory
Use `.worktrees/` for git worktrees.

## Notes
- All auth routes have stricter rate limiting via `middleware.AuthRateLimit()`
- Protected routes enforce subscription tier via `middleware.SubscriptionMiddleware`
- Async tasks (resume gen, auto-apply, scraping) run via Asynq (Go) and Celery (Python)
- Sensitive data encrypted with AES-256-GCM before storage
- CI runs linting, security scanning, vulnerability checks, tests, and Docker build
