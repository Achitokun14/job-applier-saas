# Auto Job Applier SaaS

AI-powered automated job application platform. Apply to hundreds of jobs with tailored resumes and cover letters.

## Overview

Auto Job Applier SaaS is a comprehensive platform that automates the job application process using AI. Built on top of [Jobs_Applier_AI_Agent_AIHawk](https://github.com/feder-cr/Jobs_Applier_AI_Agent_AIHawk), it provides a web dashboard, CLI tools, and a TUI for managing job applications at scale.

## Features

- **Smart Resume Generation** - AI-tailored resumes for each job application
- **Cover Letter Creation** - Personalized cover letters based on job descriptions
- **Multi-Platform Search** - Search jobs from LinkedIn, Indeed, Glassdoor, and more
- **Application Tracking** - Track all applications with status updates
- **Job Description Parsing** - Extract structured data from job postings
- **LLM Integration** - Support for OpenAI, Anthropic, Google Gemini, and Ollama
- **TUI Interface** - Terminal-based interface for local testing
- **Docker Ready** - Full docker-compose setup for deployment

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Caddy (Reverse Proxy)               │
│                         Port 80/443                         │
└───────────────┬─────────────────────────┬───────────────────┘
                │                         │
        ┌───────▼───────┐         ┌───────▼───────┐
        │   Frontend    │         │    Backend    │
        │  SvelteKit    │◄────────│   Go + Chi    │
        │  Port 3000    │         │   Port 8080   │
        └───────────────┘         └───────┬───────┘
                                          │
                          ┌───────────────┼───────────────┐
                          │               │               │
                  ┌───────▼───────┐ ┌─────▼─────┐ ┌──────▼──────┐
                  │ Python Service│ │ PostgreSQL│ │   SQLite    │
                  │   FastAPI     │ │ (Prod)    │ │   (Dev)     │
                  │  Port 8001    │ └───────────┘ └─────────────┘
                  └───────────────┘
                          ▲
                          │
                  ┌───────┴───────┐
                  │     TUI       │
                  │  Bubble Tea   │
                  └───────────────┘
```

## Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Frontend** | SvelteKit + Bun | Web dashboard UI |
| **Backend API** | Go + Chi | REST API server |
| **Resume Service** | Python + FastAPI | AI resume/cover letter generation |
| **TUI** | Go + Bubble Tea | Terminal interface for local testing |
| **Database** | SQLite (dev) / PostgreSQL (prod) | Data persistence |
| **Reverse Proxy** | Caddy | HTTPS termination, routing |
| **Containerization** | Docker + Docker Compose | Service orchestration |

## Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd job-applier-saas

# Copy environment file
cp .env.example .env

# Start all services with Docker
docker compose up -d

# Or start services individually (see GET_STARTED.md)
```

## Documentation

| Document | Description |
|----------|-------------|
| [GET_STARTED.md](GET_STARTED.md) | Installation and setup guide |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System architecture details |
| [APIs.md](APIs.md) | Complete API reference |
| [COMMANDS.md](COMMANDS.md) | All CLI commands |
| [CHANGELOG.md](CHANGELOG.md) | Version history |
| [SECURITY.md](SECURITY.md) | Security policies |
| [LICENSE.md](LICENSE.md) | License information |

### Technical Documentation

| Component | Documentation |
|-----------|---------------|
| Backend | [docs/backend/TECHNICAL.md](docs/backend/TECHNICAL.md) |
| Frontend | [docs/frontend/TECHNICAL.md](docs/frontend/TECHNICAL.md) |
| Python Service | [docs/python-service/TECHNICAL.md](docs/python-service/TECHNICAL.md) |
| TUI | [docs/tui/TECHNICAL.md](docs/tui/TECHNICAL.md) |

## Development

```bash
# Backend
cd backend && go run cmd/server/main.go

# Frontend
cd frontend && npm install && npm run dev

# Python Service
cd python-service && pip install -r requirements.txt && uvicorn main:app --port 8001

# TUI
cd tui && go run cmd/main.go
```

## Deployment

```bash
# Production deployment with Docker
docker compose -f docker-compose.production.yml up -d

# See GET_STARTED.md for detailed deployment instructions
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting pull requests.

## License

This project is proprietary software owned by Achitokun14. See [LICENSE.md](LICENSE.md) for details.

The resume generation functionality is based on [Jobs_Applier_AI_Agent_AIHawk](https://github.com/feder-cr/Jobs_Applier_AI_Agent_AIHawk) which is licensed under AGPL-3.0.

## Support

- GitHub Issues: [Create an issue](../../issues)
- Documentation: See the `/docs` directory

## Acknowledgments

- [Jobs_Applier_AI_Agent_AIHawk](https://github.com/feder-cr/Jobs_Applier_AI_Agent_AIHawk) - Original resume generation logic
- [SvelteKit](https://kit.svelte.dev/) - Frontend framework
- [Chi](https://github.com/go-chi/chi) - Go HTTP router
- [FastAPI](https://fastapi.tiangolo.com/) - Python web framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
