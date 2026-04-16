# Commands

Complete reference for all CLI commands available in Auto Job Applier SaaS.

## Backend Commands

### Start Server

```bash
# Development mode
cd backend
go run cmd/server/main.go

# Build and run
cd backend
go build -o bin/server cmd/server/main.go
./bin/server
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Server listening port |
| `DATABASE_URL` | `sqlite:./data/jobapplier.db` | Database connection string |
| `JWT_SECRET` | `default-secret-change-in-production` | JWT signing secret |
| `JWT_EXPIRY` | `24h` | JWT token expiry duration |
| `PYTHON_SERVICE_URL` | `http://localhost:8001` | Python service URL |
| `GIN_MODE` | `debug` | Gin mode (debug/release) |

### Run Tests

```bash
cd backend
go test ./...
go test -v ./...
go test -cover ./...
```

### Build

```bash
cd backend
go build -o bin/server cmd/server/main.go
```

## Frontend Commands

### Development Server

```bash
cd frontend
bun run dev
# or
npm run dev
```

### Build for Production

```bash
cd frontend
bun run build
# or
npm run build
```

### Preview Production Build

```bash
cd frontend
bun run preview
# or
npm run preview
```

### Type Checking

```bash
cd frontend
bun run check
# or
npm run check
```

### Linting

```bash
cd frontend
bun run lint
# or
npm run lint
```

### Format Code

```bash
cd frontend
bun run format
# or
npm run format
```

## Python Service Commands

### Start Server

```bash
# Development mode (with auto-reload)
cd python-service
uvicorn main:app --reload --port 8001

# Production mode
cd python-service
uvicorn main:app --host 0.0.0.0 --port 8001
```

### Install Dependencies

```bash
cd python-service
pip install -r requirements.txt
```

### Run Tests

```bash
cd python-service
python -m pytest tests/ -v
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LLM_API_KEY` | - | LLM API key (OpenAI, Anthropic, etc.) |
| `LLM_MODEL` | `gpt-4o-mini` | LLM model to use |
| `CHROME_BIN` | - | Path to Chrome binary |
| `CHROMEDRIVER_PATH` | - | Path to ChromeDriver |

## TUI Commands

### Start TUI

```bash
cd tui
go run cmd/main.go
```

### Build TUI

```bash
cd tui
go build -o bin/tui cmd/main.go
./bin/tui
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKEND_URL` | `http://localhost:8080` | Backend API URL |

## Docker Commands

### Start All Services

```bash
# Development
docker compose up -d

# Production
docker compose -f docker-compose.production.yml up -d
```

### Stop All Services

```bash
docker compose down
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f python-service
```

### Rebuild Images

```bash
docker compose build
docker compose build --no-cache
```

### Restart Services

```bash
docker compose restart
docker compose restart backend
```

### Check Status

```bash
docker compose ps
```

### Execute Commands in Container

```bash
docker compose exec backend sh
docker compose exec frontend sh
docker compose exec python-service bash
```

## Git Commands

### Initialize Repository

```bash
git init
git add .
git commit -m "Initial commit"
```

### Create Branch

```bash
git checkout -b feature/new-feature
```

### Push Changes

```bash
git add .
git commit -m "Add new feature"
git push origin feature/new-feature
```

## Database Commands

### SQLite

```bash
# Create database directory
mkdir -p backend/data

# Connect to database
sqlite3 backend/data/jobapplier.db

# List tables
.tables

# Show schema
.schema

# Exit
.exit
```

### PostgreSQL (Production)

```bash
# Connect
psql postgresql://jobapplier:password@localhost:5432/jobapplier

# List tables
\dt

# Describe table
\d users

# Exit
\q
```

## API Testing Commands

### Health Check

```bash
curl http://localhost:8080/health
```

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Search Jobs (with auth)

```bash
curl http://localhost:8080/api/v1/jobs?q=software+engineer \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Profile

```bash
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Settings

```bash
curl -X PUT http://localhost:8080/api/v1/settings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"llm_provider":"openai","llm_model":"gpt-4o-mini"}'
```

## Utility Commands

### Generate JWT Secret

```bash
openssl rand -base64 64
```

### Generate Database Password

```bash
openssl rand -base64 16
```

### Check Port Usage

```bash
# Linux/Mac
lsof -i :8080
netstat -tulpn | grep 8080

# Windows
netstat -ano | findstr :8080
```

### Kill Process on Port

```bash
# Linux/Mac
kill $(lsof -t -i:8080)

# Windows
taskkill /PID <PID> /F
```

## CI/CD Commands

### Run CI Locally

```bash
# Install act (GitHub Actions local runner)
# https://github.com/nektos/act

act -j backend
act -j frontend
act -j python-service
```

### Lint All

```bash
# Backend
cd backend && golangci-lint run

# Frontend
cd frontend && bun run lint

# Python
cd python-service && ruff check .
```
