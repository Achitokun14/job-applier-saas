@echo off
echo Starting Python Service...
cd /d "%~dp0..\python-service"
uvicorn main:app --reload --port 8001
