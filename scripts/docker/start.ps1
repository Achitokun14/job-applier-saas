@echo off
echo Starting Docker Compose...
cd /d "%~dp0.."
docker compose up -d
echo Services started. Use 'docker compose logs -f' to view logs.
