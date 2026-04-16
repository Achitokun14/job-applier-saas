@echo off
echo Running Backend Tests...
cd /d "%~dp0..\backend"
go test ./... -v
