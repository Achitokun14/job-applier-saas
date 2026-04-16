#!/bin/bash
set -e

echo "============================================"
echo "  Installing All Dependencies"
echo "============================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo ""
echo "[1/3] Installing Backend Dependencies..."
bash "$SCRIPT_DIR/backend/install.sh"

echo ""
echo "[2/3] Installing Frontend Dependencies..."
bash "$SCRIPT_DIR/frontend/install.sh"

echo ""
echo "[3/3] Installing Python Dependencies..."
bash "$SCRIPT_DIR/python-service/install.sh"

echo ""
echo "All dependencies installed!"
