@echo off
echo Viewing Docker Logs...
cd /d "%~dp0.."
docker compose logs -f
