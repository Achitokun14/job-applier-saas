#!/bin/bash
set -e

echo "============================================"
echo "  Stopping All Services"
echo "============================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Stop Docker services if running
if docker compose ps -q 2>/dev/null | grep -q .; then
    echo "Stopping Docker services..."
    docker compose down
fi

# Kill any local processes on known ports
echo "Stopping local processes..."
for PORT in 8080 5173 8001; do
    PID=$(lsof -ti:$PORT 2>/dev/null || true)
    if [ -n "$PID" ]; then
        echo "Stopping process on port $PORT (PID: $PID)"
        kill $PID 2>/dev/null || true
    fi
done

echo "All services stopped."
