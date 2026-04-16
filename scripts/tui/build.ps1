@echo off
echo Building TUI...
cd /d "%~dp0..\tui"
go build -o bin\tui cmd\main.go
echo Build complete: tui\bin\tui
