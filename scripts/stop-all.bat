@echo off
echo ============================================
echo   Stopping All Services
echo ============================================

cd /d "%~dp0.."

echo Stopping Docker services...
docker compose down 2>nul

echo Stopping local processes...
for %%P in (8080 5173 8001) do (
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :%%P ^| findstr LISTENING 2^>nul') do (
        echo Stopping process on port %%P (PID: %%a)
        taskkill /PID %%a /F >nul 2>&1
    )
)

echo All services stopped.
pause
