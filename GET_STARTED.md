# Get Started

This guide will help you set up and run Auto Job Applier SaaS on your local machine.

## Prerequisites

### Required Software

| Software | Version | Installation |
|----------|---------|--------------|
| Go | 1.22+ | [golang.org/dl](https://golang.org/dl) |
| Bun | 1.0+ | [bun.sh](https://bun.sh) |
| Python | 3.11+ | [python.org](https://python.org) |
| Git | 2.0+ | [git-scm.com](https://git-scm.com) |

### Optional Software

| Software | Purpose | Installation |
|----------|---------|--------------|
| Docker | Containerized deployment | [docker.com](https://docker.com) |
| PostgreSQL | Production database | [postgresql.org](https://postgresql.org) |
| Chrome/Chromium | PDF generation | [google.com/chrome](https://google.com/chrome) |

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd job-applier-saas
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
# Database (SQLite for local development)
DATABASE_URL=sqlite:./data/jobapplier.db

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRY=24h

# Python Service URL
PYTHON_SERVICE_URL=http://localhost:8001

# LLM Configuration (optional - for AI features)
LLM_API_KEY=your-openai-api-key
LLM_MODEL=gpt-4o-mini
```

### 3. Install Dependencies

#### Backend (Go)

```bash
cd backend
go mod download
cd ..
```

#### Frontend (Bun)

```bash
cd frontend
bun install
cd ..
```

#### Python Service

```bash
cd python-service
pip install -r requirements.txt
cd ..
```

#### TUI (Go)

```bash
cd tui
go mod download
cd ..
```

### 4. Start Services

#### Option A: Start All Services (Manual)

Open separate terminal windows for each service:

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/server/main.go
```
Backend will start on http://localhost:8080

**Terminal 2 - Frontend:**
```bash
cd frontend
bun run dev
```
Frontend will start on http://localhost:5173

**Terminal 3 - Python Service:**
```bash
cd python-service
uvicorn main:app --reload --port 8001
```
Python service will start on http://localhost:8001

**Terminal 4 - TUI (optional):**
```bash
cd tui
go run cmd/main.go
```

#### Option B: Docker Compose

```bash
# Development
docker compose up -d

# Production
docker compose -f docker-compose.production.yml up -d
```

## Verification

### 1. Check Backend

```bash
curl http://localhost:8080/health
# Expected: OK
```

### 2. Check Frontend

Open http://localhost:5173 in your browser.

### 3. Check Python Service

```bash
curl http://localhost:8001/health
# Expected: {"status":"healthy","service":"resume-generator"}
```

### 4. Test Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

## First Steps

### 1. Register an Account

1. Open http://localhost:5173
2. Click "Register"
3. Fill in your details
4. Click "Create Account"

### 2. Configure Settings

1. Navigate to "Settings"
2. Enter your LLM API key (if using AI features)
3. Configure job search preferences
4. Click "Save Settings"

### 3. Update Profile

1. Navigate to "Profile"
2. Fill in your personal information
3. Add your experience, education, skills
4. Click "Save Profile"

### 4. Search for Jobs

1. Navigate to "Jobs"
2. Enter a search query (e.g., "Software Engineer Remote")
3. Click "Search"
4. Review results and apply

### 5. Track Applications

1. Navigate to "Applications"
2. View your application history
3. Track status updates

## TUI Usage

The TUI provides a terminal-based interface for local testing:

```bash
cd tui
go run cmd/main.go
```

### TUI Controls

| Key | Action |
|-----|--------|
| Tab | Navigate between fields |
| Enter | Submit/Select |
| Esc | Go back |
| Ctrl+C | Quit |
| 1-6 | Quick menu selection |

## Troubleshooting

### Backend Won't Start

**Error:** `Failed to connect to database`

**Solution:**
```bash
mkdir -p backend/data
# Check DATABASE_URL in .env
```

### Frontend Won't Start

**Error:** `Cannot find module`

**Solution:**
```bash
cd frontend
rm -rf node_modules
bun install
```

### Python Service Errors

**Error:** `ModuleNotFoundError`

**Solution:**
```bash
cd python-service
pip install -r requirements.txt
```

**Error:** `Chrome not found`

**Solution:**
Install Chrome/Chromium or set `CHROME_BIN` environment variable.

### TUI Connection Issues

**Error:** `Connection refused`

**Solution:**
Ensure backend is running on port 8080.

## Next Steps

- Read [ARCHITECTURE.md](ARCHITECTURE.md) for system design details
- See [APIs.md](APIs.md) for complete API reference
- Check [COMMANDS.md](COMMANDS.md) for all CLI commands
- Review [SECURITY.md](SECURITY.md) for security best practices

## Support

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section
2. Search existing GitHub issues
3. Create a new issue with:
   - Error message
   - Steps to reproduce
   - Environment details
