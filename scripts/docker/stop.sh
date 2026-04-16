#!/bin/bash
set -e

echo "Stopping Docker Compose..."
cd "$(dirname "$0")/.."
docker compose down
echo "Services stopped."
