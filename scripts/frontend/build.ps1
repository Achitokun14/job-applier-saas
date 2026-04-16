@echo off
echo Building Frontend...
cd /d "%~dp0..\frontend"
npm run build
echo Build complete: frontend\build\
