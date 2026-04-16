#!/bin/bash
set -euo pipefail

# Database backup script
# Usage: ./scripts/backup.sh [output_dir]

OUTPUT_DIR="${1:-./backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${OUTPUT_DIR}/jobapplier_${TIMESTAMP}.sql.gz"

mkdir -p "$OUTPUT_DIR"

# Get database URL from env or docker
DB_URL="${DATABASE_URL:-postgres://jobapplier:changeme@localhost:5432/jobapplier}"

echo "Backing up database to ${BACKUP_FILE}..."

if command -v pg_dump &> /dev/null; then
    pg_dump "$DB_URL" | gzip > "$BACKUP_FILE"
elif docker ps | grep -q postgres; then
    docker compose exec -T postgres pg_dump -U jobapplier jobapplier | gzip > "$BACKUP_FILE"
else
    echo "Error: pg_dump not found and no running postgres container"
    exit 1
fi

echo "Backup complete: ${BACKUP_FILE} ($(du -h "$BACKUP_FILE" | cut -f1))"

# Cleanup backups older than 30 days
find "$OUTPUT_DIR" -name "jobapplier_*.sql.gz" -mtime +30 -delete 2>/dev/null || true
echo "Cleaned up backups older than 30 days"
