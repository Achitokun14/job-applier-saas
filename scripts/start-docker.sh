#!/bin/bash
set -e

echo "============================================"
echo "  Starting All Services (Docker)"
echo "============================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Check if .env exists
if [ ! -f .env ]; then
    echo "Warning: .env file not found. Using .env.example"
    cp .env.example .env
fi

echo "Building and starting containers..."
docker compose up -d --build

echo ""
echo "All services started!"
echo "  - Frontend:     http://localhost:3000"
echo "  - Backend API:  http://localhost:8080"
echo "  - Python:       http://localhost:8001"
echo ""
echo "View logs: docker compose logs -f"
echo "Stop:      docker compose down"
