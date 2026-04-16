@echo off
chcp 65001 >nul 2>&1
setlocal EnableDelayedExpansion

:: ============================================================
::  Job Applier SaaS - Main Menu (Windows Batch)
:: ============================================================

set "SCRIPT_DIR=%~dp0"

:menu
cls
echo ============================================================
echo          Job Applier SaaS - Main Menu
echo          Automated Job Application SaaS
echo ============================================================
echo.
echo   === Run All Services ===
echo   1^) Start All (Local)
echo   2^) Start All (Docker)
echo   3^) Stop All Services
echo.
echo   === Individual Services ===
echo   4^) Start Backend Only
echo   5^) Start Frontend Only
echo   6^) Start Python Service Only
echo   7^) Start TUI Only
echo.
echo   === Build and Install ===
echo   8^) Install All Dependencies
echo   9^) Build All
echo   10^) Build Backend
echo   11^) Build Frontend
echo   12^) Build TUI
echo.
echo   === Docker ===
echo   13^) Docker Build
echo   14^) Docker Logs
echo   15^) Docker Stop
echo.
echo   === Utilities ===
echo   0^) Exit
echo.
set /p choice="Select option: "

if "%choice%"=="1" goto start_all
if "%choice%"=="2" goto start_docker
if "%choice%"=="3" goto stop_all
if "%choice%"=="4" goto start_backend
if "%choice%"=="5" goto start_frontend
if "%choice%"=="6" goto start_python
if "%choice%"=="7" goto start_tui
if "%choice%"=="8" goto install_all
if "%choice%"=="9" goto build_all
if "%choice%"=="10" goto build_backend
if "%choice%"=="11" goto build_frontend
if "%choice%"=="12" goto build_tui
if "%choice%"=="13" goto docker_build
if "%choice%"=="14" goto docker_logs
if "%choice%"=="15" goto docker_stop
if "%choice%"=="0" goto exit

echo Invalid option.
pause
goto menu

:start_all
call "%SCRIPT_DIR%scripts\start-all.bat"
pause
goto menu

:start_docker
call "%SCRIPT_DIR%scripts\start-docker.bat"
pause
goto menu

:stop_all
call "%SCRIPT_DIR%scripts\stop-all.bat"
pause
goto menu

:start_backend
call "%SCRIPT_DIR%scripts\backend\start.bat"
pause
goto menu

:start_frontend
call "%SCRIPT_DIR%scripts\frontend\start.bat"
pause
goto menu

:start_python
call "%SCRIPT_DIR%scripts\python-service\start.bat"
pause
goto menu

:start_tui
call "%SCRIPT_DIR%scripts\tui\start.bat"
pause
goto menu

:install_all
call "%SCRIPT_DIR%scripts\install-all.bat"
pause
goto menu

:build_all
echo Building all services...
call "%SCRIPT_DIR%scripts\backend\build.bat"
call "%SCRIPT_DIR%scripts\frontend\build.bat"
call "%SCRIPT_DIR%scripts\tui\build.bat"
echo All builds complete!
pause
goto menu

:build_backend
call "%SCRIPT_DIR%scripts\backend\build.bat"
pause
goto menu

:build_frontend
call "%SCRIPT_DIR%scripts\frontend\build.bat"
pause
goto menu

:build_tui
call "%SCRIPT_DIR%scripts\tui\build.bat"
pause
goto menu

:docker_build
call "%SCRIPT_DIR%scripts\docker\build.bat"
pause
goto menu

:docker_logs
call "%SCRIPT_DIR%scripts\docker\logs.bat"
pause
goto menu

:docker_stop
call "%SCRIPT_DIR%scripts\docker\stop.bat"
pause
goto menu

:exit
echo Goodbye!
exit /b 0
