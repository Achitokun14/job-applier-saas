#!/bin/bash
set -e

echo "Building Backend..."
cd "$(dirname "$0")/../backend"
go build -o bin/server cmd/server/main.go
echo "Build complete: backend/bin/server"
