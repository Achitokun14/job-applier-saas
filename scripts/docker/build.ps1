@echo off
echo Building Docker Images...
cd /d "%~dp0.."
docker compose build
echo Build complete.
