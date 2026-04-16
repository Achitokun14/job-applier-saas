@echo off
echo Running Frontend Type Check...
cd /d "%~dp0..\frontend"
npm run check
