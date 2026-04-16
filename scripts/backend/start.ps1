@echo off
echo Starting Backend Server...
cd /d "%~dp0..\backend"
go run cmd\server\main.go
