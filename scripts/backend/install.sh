#!/bin/bash
set -e

echo "Installing Backend Dependencies..."
cd "$(dirname "$0")/../backend"
go mod download
echo "Dependencies installed."
