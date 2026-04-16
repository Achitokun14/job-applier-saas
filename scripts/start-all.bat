@echo off
echo ============================================
echo   Starting All Services (Local Development)
echo ============================================

set SCRIPT_DIR=%~dp0

echo Starting Backend...
start "Backend" cmd /c "%SCRIPT_DIR%backend\start.bat"
timeout /t 2 /nobreak >nul

echo Starting Frontend...
start "Frontend" cmd /c "%SCRIPT_DIR%frontend\start.bat"

echo Starting Python Service...
start "Python Service" cmd /c "%SCRIPT_DIR%python-service\start.bat"

echo.
echo All services started!
echo   - Backend:      http://localhost:8080
echo   - Frontend:     http://localhost:5173
echo   - Python:       http://localhost:8001
echo.
echo Press any key to stop all services...
pause >nul

echo Stopping services...
taskkill /FI "WINDOWTITLE eq Backend*" /F >nul 2>&1
taskkill /FI "WINDOWTITLE eq Frontend*" /F >nul 2>&1
taskkill /FI "WINDOWTITLE eq Python Service*" /F >nul 2>&1
echo Done.
