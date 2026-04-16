#!/bin/bash
set -e

echo "============================================"
echo "  Job Applier SaaS - Update Script"
echo "============================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Pull latest code
echo "Pulling latest code..."
git pull origin main

# Rebuild and restart
echo "Rebuilding services..."
docker compose -f docker-compose.production.yml build --no-cache

echo "Restarting services..."
docker compose -f docker-compose.production.yml down
docker compose -f docker-compose.production.yml up -d

echo ""
echo "Update complete!"
echo "Check status: docker compose -f docker-compose.production.yml ps"
