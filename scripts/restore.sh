#!/bin/bash
set -euo pipefail

# Database restore script
# Usage: ./scripts/restore.sh <backup_file>

BACKUP_FILE="${1:?Usage: ./scripts/restore.sh <backup_file.sql.gz>}"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: File not found: $BACKUP_FILE"
    exit 1
fi

DB_URL="${DATABASE_URL:-postgres://jobapplier:changeme@localhost:5432/jobapplier}"

echo "WARNING: This will overwrite the current database!"
read -p "Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

echo "Restoring from ${BACKUP_FILE}..."
gunzip -c "$BACKUP_FILE" | psql "$DB_URL"
echo "Restore complete."
