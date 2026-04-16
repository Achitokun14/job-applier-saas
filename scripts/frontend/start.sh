#!/bin/bash
set -e

echo "Starting Frontend Dev Server..."
cd "$(dirname "$0")/../frontend"
npm run dev
