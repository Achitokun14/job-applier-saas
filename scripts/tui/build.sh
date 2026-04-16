#!/bin/bash
set -e

echo "Building TUI..."
cd "$(dirname "$0")/../tui"
go build -o bin/tui cmd/main.go
echo "Build complete: tui/bin/tui"
