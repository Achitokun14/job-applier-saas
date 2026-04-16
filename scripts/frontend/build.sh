#!/bin/bash
set -e

echo "Building Frontend..."
cd "$(dirname "$0")/../frontend"
npm run build
echo "Build complete: frontend/build/"
