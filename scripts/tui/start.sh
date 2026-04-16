#!/bin/bash
set -e

echo "Starting TUI..."
cd "$(dirname "$0")/../tui"
go run cmd/main.go
