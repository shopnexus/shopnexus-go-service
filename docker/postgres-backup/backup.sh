#!/bin/bash
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="shopnexus"

mkdir -p "$BACKUP_DIR"

echo "$(date): Starting backup of $DB_NAME"

pg_dump -Fc $DB_NAME >$BACKUP_DIR/backup_${DB_NAME}_${DATE}.dump

if [ -f "$BACKUP_DIR/backup_${DB_NAME}_${DATE}.dump" ]; then
  echo "$(date): Backup created successfully: backup_${DB_NAME}_${DATE}.dump"
  find "$BACKUP_DIR" -name "backup_${DB_NAME}_*.dump" -mtime +7 -delete
  echo "$(date): Old backups cleaned up"
else
  echo "$(date): ERROR: Backup failed!"
  exit 1
fi
