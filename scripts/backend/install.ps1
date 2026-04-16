@echo off
echo Installing Backend Dependencies...
cd /d "%~dp0..\backend"
go mod download
echo Dependencies installed.
