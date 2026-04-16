#!/bin/bash
set -e

echo "Running Backend Tests..."
cd "$(dirname "$0")/../backend"
go test ./... -v
