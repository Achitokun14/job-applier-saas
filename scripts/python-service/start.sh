#!/bin/bash
set -e

echo "Starting Python Service..."
cd "$(dirname "$0")/../python-service"
uvicorn main:app --reload --port 8001
