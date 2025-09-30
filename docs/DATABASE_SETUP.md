# Database Setup Guide

## 🚀 Quick Start

### Step 1: Start PostgreSQL

You have **two options**:

#### Option A: Using Docker (Recommended for Consistency)

```bash
docker-compose up -d
```

This starts PostgreSQL in a container with the settings from `docker-compose.yml`.

**Check if it's running:**
```bash
docker ps
# You should see: go-ddd-postgres
```

#### Option B: Using Local PostgreSQL (If Already Installed)

If you have PostgreSQL installed locally, just make sure it's running:

**Windows:**
- PostgreSQL should be running as a service
- Check: Services → PostgreSQL

**Linux/Mac:**
```bash
# Check if running
sudo systemctl status postgresql
# or
pg_isready

# Start if not running
sudo systemctl start postgresql
```

**Update `.env` to match your local setup:**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres          # Your PostgreSQL username
DB_PASSWORD=your_password # Your PostgreSQL password
DB_NAME=go_ddd_starter
```

### Step 2: Create Database

```bash
bash scripts/migrate.sh create
```

This creates the `go_ddd_starter` database.

### Step 3: Run Migrations

```bash
bash scripts/migrate.sh up
```

This creates all tables, indexes, and schema.

### Step 4: Verify

```bash
# Connect to database
psql -U postgres -h localhost -d go_ddd_starter

# Inside psql:
\dt              # List tables
\d users         # Describe users table
\q               # Quit
```

---

## 📝 Understanding the Migration System

### Migration Files

Migrations are numbered SQL files:
```
000001_init_schema.up.sql       # First migration (creates initial schema)
000001_init_schema.down.sql     # Rollback for first migration
000002_add_user_roles.up.sql    # Second migration (adds features)
000002_add_user_roles.down.sql  # Rollback for second migration
```

**Naming pattern:** `{number}_{description}.{direction}.sql`
- **number**: Sequential (000001, 000002, ...)
- **description**: What this migration does
- **direction**: `up` (apply) or `down` (rollback)

### How Migrations Work

**UP migrations** (`*.up.sql`):
- Apply changes to database
- Create tables, add columns, create indexes
- Run in numerical order (000001 → 000002 → 000003)

**DOWN migrations** (`*.down.sql`):
- Undo changes from UP migration
- Drop tables, remove columns, drop indexes
- Run in reverse order (000003 → 000002 → 000001)

### The Migration Script

`scripts/migrate.sh` is a simple helper that:
1. Loads database credentials from `.env`
2. Runs SQL files using `psql`
3. Provides commands: `create`, `up`, `down`, `drop`, `reset`

**Commands:**
```bash
bash scripts/migrate.sh create  # Create database
bash scripts/migrate.sh up      # Run all UP migrations
bash scripts/migrate.sh down    # Run all DOWN migrations (rollback)
bash scripts/migrate.sh drop    # Delete database (asks for confirmation)
bash scripts/migrate.sh reset   # Drop + Create + Migrate (fresh start)
```

---

## 🔧 How to Add New Migrations

### Example: Adding User Roles

**1. Create migration files:**
```bash
# Create empty files
touch infrastructure/database/migrations/000002_add_user_roles.up.sql
touch infrastructure/database/migrations/000002_add_user_roles.down.sql
```

**2. Write UP migration** (`000002_add_user_roles.up.sql`):
```sql
-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user';

-- Create index for role queries
CREATE INDEX idx_users_role ON users(role);

-- Add comment
COMMENT ON COLUMN users.role IS 'User role: admin, user, guest';
```

**3. Write DOWN migration** (`000002_add_user_roles.down.sql`):
```sql
-- Remove index
DROP INDEX IF EXISTS idx_users_role;

-- Remove column
ALTER TABLE users DROP COLUMN IF EXISTS role;
```

**4. Run the migration:**
```bash
bash scripts/migrate.sh up
```

**5. Test rollback:**
```bash
bash scripts/migrate.sh down  # Undo
bash scripts/migrate.sh up    # Reapply
```

---

## 🗄️ Current Schema (After Migration 000001)

**Users Table:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Indexes:**
- `idx_users_email` - Fast email lookups (for login)
- `idx_users_is_active` - Filter active/deleted users
- `idx_users_created_at` - Sort by creation date

---

## 🔍 SQLC Integration

After migrations are applied, SQLC uses the schema to generate Go code.

**How it works:**
1. You write SQL queries in `internal/domains/users/infrastructure/persistence/queries/users.sql`
2. SQLC reads the schema from migrations
3. SQLC generates type-safe Go code in `internal/domains/users/infrastructure/persistence/sqlc/`

**Generate SQLC code:**
```bash
make sqlc-generate
# or
sqlc generate
```

---

## ⚙️ Configuration

### Database Connection (`.env`)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_ddd_starter
DB_SSLMODE=disable
```

### Docker Compose (`docker-compose.yml`)

```yaml
services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: go_ddd_starter
```

---

## 🛠️ Troubleshooting

### "Database does not exist"
```bash
bash scripts/migrate.sh create
```

### "Connection refused"
```bash
# Check if PostgreSQL is running
docker ps

# Start it if not running
docker-compose up -d
```

### "Permission denied" on script
```bash
chmod +x scripts/migrate.sh
```

### Reset everything
```bash
bash scripts/migrate.sh reset
```

This will:
1. Drop the database
2. Create a fresh database
3. Run all migrations

---

## 📌 Best Practices

1. **Never modify existing migrations** - Create new ones instead
2. **Always test rollback** - Run `down` then `up` to verify
3. **Use transactions** - Wrap DDL in `BEGIN`/`COMMIT` when possible
4. **Add comments** - Document what columns/tables are for
5. **Create indexes** - For columns you'll query frequently
6. **Backup before production** - Always backup before migrating in production

---

## 🎯 Summary

**To get started:**
```bash
# 1. Start PostgreSQL
docker-compose up -d

# 2. Create database and run migrations
bash scripts/migrate.sh create
bash scripts/migrate.sh up

# 3. Verify
psql -U postgres -h localhost -d go_ddd_starter -c "\dt"
```

**To add new features:**
```bash
# 1. Create migration files (000002, 000003, etc.)
# 2. Write SQL for up and down
# 3. Run migration
bash scripts/migrate.sh up
```

That's it! The migration system is simple but powerful. 🚀
