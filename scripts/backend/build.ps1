@echo off
echo Building Backend...
cd /d "%~dp0..\backend"
go build -o bin\server cmd\server\main.go
echo Build complete: backend\bin\server
