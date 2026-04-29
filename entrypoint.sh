#!/bin/sh
# entrypoint.sh
# Runs SQL migrations in order then starts the Go server.
# Using psql (from the postgres client package) keeps the migration
# runner simple — no extra Go dependencies or migration framework needed.

set -e

echo "==> Running migrations..."

# Parse host/port/user/password/dbname from the DSN in config.yaml.
# We rely on the PGPASSWORD env var (set in docker-compose) so psql
# doesn't prompt for a password.
PGHOST="${PGHOST:-db}"
PGPORT="${PGPORT:-5432}"
PGUSER="${PGUSER:-kfc}"
PGDATABASE="${PGDATABASE:-kfc}"

# Wait until Postgres is ready (it can take a few seconds to init).
until pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -q; do
  echo "   waiting for postgres..."
  sleep 1
done

# Run every .sql file in order. Files are named 001_, 002_, etc.
for f in ./migrations/*.sql; do
  echo "   applying $f"
  psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" -f "$f"
done

echo "==> Migrations done. Starting server..."
exec ./kfc-server
