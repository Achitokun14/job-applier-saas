#!/bin/bash
set -e

echo "Viewing Docker Logs..."
cd "$(dirname "$0")/.."
docker compose logs -f
