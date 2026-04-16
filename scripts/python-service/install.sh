#!/bin/bash
set -e

echo "Installing Python Dependencies..."
cd "$(dirname "$0")/../python-service"
pip install -r requirements.txt
echo "Dependencies installed."
