#!/bin/bash
set -e

echo "Running Frontend Type Check..."
cd "$(dirname "$0")/../frontend"
npm run check
