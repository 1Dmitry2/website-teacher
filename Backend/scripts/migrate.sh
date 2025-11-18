#!/bin/sh

# Migration script using golang-migrate
# Install golang-migrate first: https://github.com/golang-migrate/migrate

MIGRATIONS_DIR="./migrations"
DATABASE_URL="${DATABASE_URL:-host=localhost dbname=restapi_dev user=postgres password=postgres sslmode=disable}"

if [ -z "$1" ]; then
    echo "Usage: $0 [up|down|version]"
    exit 1
fi

if ! command -v migrate &> /dev/null; then
    echo "Error: golang-migrate is not installed"
    echo "Install it from: https://github.com/golang-migrate/migrate"
    exit 1
fi

case "$1" in
    up)
        migrate -path "$MIGRATIONS_DIR" -database "postgres://$DATABASE_URL" up
        ;;
    down)
        migrate -path "$MIGRATIONS_DIR" -database "postgres://$DATABASE_URL" down
        ;;
    version)
        migrate -path "$MIGRATIONS_DIR" -database "postgres://$DATABASE_URL" version
        ;;
    *)
        echo "Unknown command: $1"
        echo "Usage: $0 [up|down|version]"
        exit 1
        ;;
esac

