@echo off
echo ============================================
echo   Starting All Services (Docker)
echo ============================================

cd /d "%~dp0.."

if not exist .env (
    echo Warning: .env file not found. Using .env.example
    copy .env.example .env
)

echo Building and starting containers...
docker compose up -d --build

echo.
echo All services started!
echo   - Frontend:     http://localhost:3000
echo   - Backend API:  http://localhost:8080
echo   - Python:       http://localhost:8001
echo.
echo View logs: docker compose logs -f
echo Stop:      docker compose down
pause
