# Go DDD Clean Architecture Starter

A production-ready Go starter project implementing **Domain-Driven Design (DDD)** with **Clean Architecture** principles.

## 🏗️ Architecture

This project follows a domain-centric architecture where:
- **Each domain is a separate silo** (users, products, orders, etc.)
- **Clean Architecture is applied within each domain** (domain, application, infrastructure, handler layers)
- **Fiber handles HTTP routing only** (no business logic in framework)
- **No cross-domain dependencies** (domains communicate via events or APIs)

For detailed architecture documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## 🚀 Tech Stack

- **Web Framework**: [Fiber v2](https://gofiber.io/) - Fast HTTP routing
- **Database**: PostgreSQL
- **Database Driver**: [pgx/v5](https://github.com/jackc/pgx) - High-performance PostgreSQL driver
- **Query Builder**: [SQLC](https://sqlc.dev/) - Type-safe SQL code generation
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) (optional)

## 📁 Project Structure

```
go-ddd-clean-starter/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── domains/                       # Domain silos
│   │   └── users/                     # User domain
│   │       ├── domain/                # Business logic (entities, interfaces)
│   │       ├── application/           # Use cases & DTOs
│   │       ├── infrastructure/        # Repository implementation & SQLC
│   │       └── handler/               # HTTP handlers
│   └── platform/                      # Shared infrastructure
│       ├── database/                  # DB connection & transactions
│       ├── logger/                    # Logging
│       ├── middleware/                # HTTP middleware
│       ├── config/                    # Configuration
│       ├── errors/                    # Error handling
│       └── docs/                      # API documentation
├── infrastructure/
│   └── database/
│       └── migrations/                # Database migrations
├── docs/
│   ├── ARCHITECTURE.md                # Architecture documentation
│   └── openapi/
│       └── openapi.yaml               # OpenAPI specification
├── sqlc.yaml                          # SQLC configuration
├── docker-compose.yml                 # PostgreSQL container
├── Makefile                           # Development commands
└── .env.example                       # Environment variables template
```

## 🛠️ Prerequisites

Before you begin, ensure you have the following installed:

1. **Go 1.21+** - [Download](https://golang.org/dl/)
2. **PostgreSQL** - [Download](https://www.postgresql.org/download/) or use Docker
3. **SQLC** - [Installation Guide](https://docs.sqlc.dev/en/latest/overview/install.html)
   ```bash
   # macOS
   brew install sqlc
   
   # Linux/Windows - Download from releases
   # Or use: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```
4. **golang-migrate** (optional) - [Installation Guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
5. **Docker** (optional) - For running PostgreSQL in a container

## 🚦 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/dzikrisyairozi/go-ddd-clean-starter.git
cd go-ddd-clean-starter
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 4. Start PostgreSQL

**Option A: Using Docker (Recommended)**
```bash
docker-compose up -d
```

**Option B: Using Local PostgreSQL**
```bash
# Create database
createdb go_ddd_starter
```

### 5. Run Database Migrations

**Option A: Using golang-migrate**
```bash
make migrate-up
```

**Option B: Manual SQL execution**
```bash
psql -U postgres -d go_ddd_starter -f infrastructure/database/migrations/000001_create_users_table.up.sql
```

### 6. Generate SQLC Code

```bash
make sqlc-generate
# Or without sqlc installed:
make sqlc-generate-go
```

### 7. Run the Application

```bash
make run
# Or:
go run cmd/api/main.go
```

The server will start on `http://localhost:3000`

## 📚 API Documentation

Once the server is running, access the interactive API documentation:

- **Scalar UI**: http://localhost:3000/docs
- **OpenAPI Spec**: http://localhost:3000/openapi.yaml

## 🧪 Testing the API

### Create a User

```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "password": "SecurePass123!"
  }'
```

### Get User by ID

```bash
curl http://localhost:3000/api/v1/users/{user-id}
```

### Update User

```bash
curl -X PUT http://localhost:3000/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:3000/api/v1/users/{user-id}
```

## 🔧 Development Commands

```bash
make help              # Show all available commands
make run               # Run the application
make build             # Build binary
make test              # Run tests
make test-coverage     # Run tests with coverage
make sqlc-generate     # Generate SQLC code
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make docker-up         # Start PostgreSQL container
make docker-down       # Stop PostgreSQL container
make fmt               # Format code
make lint              # Run linter
make clean             # Clean build artifacts
```

## 🏛️ Architecture Principles

### Domain Layer (Core)
- Pure business logic
- No external dependencies
- Defines repository interfaces
- Contains entities and value objects

### Application Layer (Use Cases)
- Orchestrates domain objects
- Defines transaction boundaries
- Contains DTOs for input/output

### Infrastructure Layer
- Implements repository interfaces
- Uses SQLC for database access
- Maps between database models and domain entities

### Handler Layer (HTTP)
- Handles HTTP requests/responses
- Validates input
- Calls application services
- No business logic

## 📖 Adding a New Domain

```bash
# 1. Create domain structure
mkdir -p internal/domains/products/{domain,application,infrastructure/persistence/{queries,sqlc},handler}

# 2. Create migration
make migrate-create name=create_products_table

# 3. Add SQLC queries
# Create internal/domains/products/infrastructure/persistence/queries/products.sql

# 4. Update sqlc.yaml with new domain configuration

# 5. Generate SQLC code
make sqlc-generate

# 6. Implement domain layers (domain → application → infrastructure → handler)

# 7. Register routes in cmd/api/main.go
```

## 🤝 Contributing

Contributions are welcome! Please read the architecture documentation before contributing.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Clean Architecture principles by Robert C. Martin
- Domain-Driven Design by Eric Evans
- Go community best practices

## 📞 Support

For questions or issues, please open an issue on GitHub.

---

**Happy Coding! 🚀**