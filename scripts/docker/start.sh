#!/bin/bash
set -e

echo "Starting Docker Compose..."
cd "$(dirname "$0")/.."
docker compose up -d
echo "Services started. Use 'docker compose logs -f' to view logs."
