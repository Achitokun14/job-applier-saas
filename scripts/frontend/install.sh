#!/bin/bash
set -e

echo "Installing Frontend Dependencies..."
cd "$(dirname "$0")/../frontend"
npm install
echo "Dependencies installed."
