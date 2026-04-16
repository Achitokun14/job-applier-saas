@echo off
echo Stopping Docker Compose...
cd /d "%~dp0.."
docker compose down
echo Services stopped.
