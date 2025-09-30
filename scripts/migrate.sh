#!/bin/bash

# Database Migration Helper Script
# This script helps run database migrations using psql

set -e

# Load environment variables from .env if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-go_ddd_starter}

MIGRATIONS_DIR="infrastructure/database/migrations"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run migration
run_migration() {
    local file=$1
    print_info "Running migration: $file"
    
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$file"
    
    if [ $? -eq 0 ]; then
        print_info "✓ Migration successful: $file"
    else
        print_error "✗ Migration failed: $file"
        exit 1
    fi
}

# Function to migrate up
migrate_up() {
    print_info "Running UP migrations..."
    
    for file in $MIGRATIONS_DIR/*.up.sql; do
        if [ -f "$file" ]; then
            run_migration "$file"
        fi
    done
    
    print_info "All migrations completed successfully!"
}

# Function to migrate down
migrate_down() {
    print_warn "Running DOWN migrations (this will rollback changes)..."
    
    # Run down migrations in reverse order
    for file in $(ls -r $MIGRATIONS_DIR/*.down.sql); do
        if [ -f "$file" ]; then
            run_migration "$file"
        fi
    done
    
    print_info "All rollbacks completed successfully!"
}

# Function to create database if it doesn't exist
create_db() {
    print_info "Creating database: $DB_NAME"
    
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"
    
    if [ $? -eq 0 ]; then
        print_info "✓ Database created: $DB_NAME"
    else
        print_warn "Database might already exist or creation failed"
    fi
}

# Function to drop database
drop_db() {
    print_warn "Dropping database: $DB_NAME"
    read -p "Are you sure? This will delete all data! (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
        print_info "✓ Database dropped: $DB_NAME"
    else
        print_info "Operation cancelled"
    fi
}

# Main script
case "$1" in
    up)
        migrate_up
        ;;
    down)
        migrate_down
        ;;
    create)
        create_db
        ;;
    drop)
        drop_db
        ;;
    reset)
        drop_db
        create_db
        migrate_up
        ;;
    *)
        echo "Usage: $0 {up|down|create|drop|reset}"
        echo ""
        echo "Commands:"
        echo "  up     - Run all pending migrations"
        echo "  down   - Rollback all migrations"
        echo "  create - Create the database"
        echo "  drop   - Drop the database (with confirmation)"
        echo "  reset  - Drop, create, and migrate (fresh start)"
        exit 1
        ;;
esac
