@echo off
echo Installing Frontend Dependencies...
cd /d "%~dp0..\frontend"
npm install
echo Dependencies installed.
