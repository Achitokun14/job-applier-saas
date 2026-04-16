@echo off
echo Starting TUI...
cd /d "%~dp0..\tui"
go run cmd\main.go
