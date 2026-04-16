#!/bin/bash
set -e

echo "============================================"
echo "  Starting All Services (Local Development)"
echo "============================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Function to run in background
run_bg() {
    echo "Starting $1..."
    "$@" &
    echo "$1 started (PID: $!)"
}

# Start all services
run_bg bash "$SCRIPT_DIR/backend/start.sh"
sleep 2
run_bg bash "$SCRIPT_DIR/frontend/start.sh"
run_bg bash "$SCRIPT_DIR/python-service/start.sh"

echo ""
echo "All services started!"
echo "  - Backend:      http://localhost:8080"
echo "  - Frontend:     http://localhost:5173"
echo "  - Python:       http://localhost:8001"
echo ""
echo "Press Ctrl+C to stop all services"
wait
