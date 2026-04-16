#!/bin/bash
set -e

echo "Building Docker Images..."
cd "$(dirname "$0")/.."
docker compose build
echo "Build complete."
