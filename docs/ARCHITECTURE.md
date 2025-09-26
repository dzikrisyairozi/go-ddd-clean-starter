# Architecture Specification

## Overview

This project implements a **Domain-Driven Design (DDD)** architecture where **each domain is a separate silo**, and within each domain, **Clean Architecture** principles are applied.

### Technology Stack
- **Web Framework**: Fiber (HTTP routing ONLY - no business logic)
- **Database**: PostgreSQL
- **Database Driver**: pgx/v5 (PostgreSQL driver and toolkit)
- **Query Builder**: SQLC (generates type-safe Go code from SQL)

### Key Architectural Decisions
1. **Domain-Centric**: Each domain (users, products, orders, etc.) is isolated in its own module
2. **Clean Architecture Inside Each Domain**: Each domain has its own layers (domain, application, infrastructure)
3. **Fiber for Routing Only**: Fiber handles HTTP routing and delegates to domain services
4. **No Cross-Domain Dependencies**: Domains communicate through events or APIs

## Architecture Structure

### High-Level Organization

```
go-ddd-clean-starter/
├── internal/
│   ├── domains/             # Each domain is a silo
│   │   ├── users/           # User domain
│   │   ├── products/        # Product domain (example)
│   │   └── orders/          # Order domain (example)
│   └── platform/            # Technical infrastructure (no business logic)
│                            # Contains: db, config, middleware, logger, errors, utils
└── cmd/
    └── api/
        └── main.go          # Application entry point
```

## Simplified Domain Structure

Each domain is **isolated** but kept **simple** for a starter template:

```
internal/domains/users/
├── domain/                  # Domain Layer (Core Business Logic)
│   ├── user.go              # User entity + value objects
│   ├── repository.go        # Repository interface (contract)
│   └── errors.go            # Domain-specific errors
│
├── application/             # Application Layer (Use Cases)
│   ├── service.go           # Application service (all use cases)
│   └── dto.go               # Input/Output DTOs
│
├── infrastructure/          # Infrastructure Layer
│   ├── persistence/
│   │   ├── queries/         # SQLC queries for this domain
│   │   │   └── users.sql
│   │   └── sqlc/            # SQLC generated code for this domain
│   └── repository.go        # Repository implementation (uses SQLC)
│
└── handler/                 # HTTP Layer
    ├── handler.go           # HTTP handlers
    ├── request.go           # Request models
    └── response.go          # Response models
```

**Key Simplifications**:
- ✅ Merged entity + value objects into single file
- ✅ Single service file instead of one per use case
- ✅ Removed unnecessary nested folders
- ✅ Handler directly in domain (no extra presentation folder)
- ✅ Database migrations in shared infrastructure
- ✅ **Domain-specific queries live with each domain**
- ✅ **Single SQLC config at project root** (generates code for all domains)

## Layer Responsibilities

### 1. Domain Layer (Core)
**Location**: `internal/domains/{domain}/domain/`

**Purpose**: Pure business logic and domain models. **NO external dependencies**.

**Contains**:
- **Entities**: Business objects with identity and lifecycle
- **Value Objects**: Immutable objects defined by attributes
- **Repository Interfaces**: Contracts for data persistence
- **Domain Services**: Business logic that doesn't fit in entities
- **Domain Events**: Events representing state changes
- **Domain Errors**: Business rule violations

**Rules**:
- ✅ Pure Go types and interfaces
- ✅ Business rules enforced here
- ❌ NO framework dependencies
- ❌ NO infrastructure dependencies
- ❌ NO application layer dependencies

### 2. Application Layer (Use Cases)
**Location**: `internal/domains/{domain}/application/`

**Purpose**: Orchestrates domain objects to fulfill use cases.

**Contains**:
- **Application Services**: Coordinates domain objects
- **Use Cases**: Specific business flows
- **DTOs**: Data transfer objects for input/output
- **Ports**: Interfaces for external services

**Rules**:
- ✅ Depends on domain layer only
- ✅ Defines transaction boundaries
- ❌ NO framework dependencies
- ❌ NO infrastructure dependencies

### 3. Infrastructure Layer (External)
**Location**: `internal/domains/{domain}/infrastructure/`

**Purpose**: Implements repository interface from domain layer.

**Contains**:
- **Repository Implementation**: Uses SQLC generated code

**Rules**:
- ✅ Implements domain repository interface
- ✅ Uses SQLC for type-safe queries
- ✅ Receives pgx pool via dependency injection

**Note**: 
- Database migrations are in shared infrastructure (`infrastructure/database/migrations/`)
- SQLC queries are **inside each domain** (`internal/domains/{domain}/infrastructure/persistence/queries/`)
- Each domain has its own SQLC generated code

### 4. Handler Layer (HTTP)
**Location**: `internal/domains/{domain}/handler/`

**Purpose**: Handles HTTP requests and responses.

**Contains**:
- **HTTP Handlers**: Process HTTP requests
- **Request/Response Models**: HTTP-specific DTOs

**Rules**:
- ✅ Depends on application layer
- ✅ Handles HTTP status codes, headers
- ✅ Validates HTTP input
- ❌ NO business logic here
- ❌ Does NOT call domain layer directly

## Dependency Flow Within Each Domain

```
┌─────────────────────────────────────┐
│     Presentation Layer (HTTP)       │
│     - Handlers                      │
│     - Request/Response              │
└──────────────┬──────────────────────┘
               │ calls
               ↓
┌─────────────────────────────────────┐
│     Application Layer               │
│     - Use Cases                     │
│     - Application Services          │
└──────────────┬──────────────────────┘
               │ uses
               ↓
┌─────────────────────────────────────┐
│     Domain Layer (CORE)             │
│     - Entities                      │
│     - Value Objects                 │
│     - Repository Interfaces         │
└─────────────────────────────────────┘
               ↑
               │ implements
┌──────────────┴──────────────────────┐
│     Infrastructure Layer            │
│     - Repository Impl (SQLC)        │
│     - External Adapters             │
└─────────────────────────────────────┘
```

## Fiber's Role: Routing ONLY

**Fiber is ONLY used for HTTP routing and middleware**. It does NOT touch business logic.

```go
// cmd/api/main.go
package main

import (
    "github.com/gofiber/fiber/v2"
    userHandler "yourproject/internal/domains/users/handler"
    productHandler "yourproject/internal/domains/products/handler"
)

func main() {
    app := fiber.New()
    
    // Global middleware
    app.Use(logger.New())
    app.Use(cors.New())
    
    // Register domain routes
    api := app.Group("/api/v1")
    userHandler.RegisterRoutes(api, userService)      // Users domain
    productHandler.RegisterRoutes(api, productService) // Products domain
    
    app.Listen(":3000")
}
```

```go
// internal/domains/users/handler/routes.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "yourproject/internal/domains/users/application"
)

func RegisterRoutes(router fiber.Router, service application.UserService) {
    handler := NewUserHandler(service)
    
    users := router.Group("/users")
    users.Post("/", handler.CreateUser)       // POST /api/v1/users
    users.Get("/:id", handler.GetUser)        // GET /api/v1/users/:id
    users.Put("/:id", handler.UpdateUser)     // PUT /api/v1/users/:id
    users.Delete("/:id", handler.DeleteUser)  // DELETE /api/v1/users/:id
}
```

```go
// internal/domains/users/handler/handler.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "yourproject/internal/domains/users/application"
)

type UserHandler struct {
    service application.UserService  // Application layer, NOT domain
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
    }
    
    // Map HTTP request to application DTO
    dto := application.CreateUserDTO{
        Email: req.Email,
        Name:  req.Name,
    }
    
    // Call application service (use case)
    result, err := h.service.CreateUser(c.Context(), dto)
    if err != nil {
        return handleError(c, err)
    }
    
    // Map application DTO to HTTP response
    response := mapToResponse(result)
    return c.Status(201).JSON(response)
}
```

## Cross-Domain Communication

Domains are **isolated** and should NOT directly depend on each other.

### Option 1: Domain Events (Recommended)
```go
// When user is created, publish event
event := domain.UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
}
eventBus.Publish(event)

// Other domains subscribe to events
// orders domain listens for UserCreatedEvent
```

### Option 2: API Calls
```go
// One domain calls another through HTTP API
// (treat it like an external service)
```

### Option 3: Shared Kernel (Optional)
If you have truly shared domain concepts (e.g., Money, Address value objects), you can create:
```go
// internal/shared/
// Only for domain concepts shared across multiple domains
// Keep this MINIMAL - most code should live in domains
```

**Note**: This is optional. Only add if you actually need it.

## Shared Infrastructure

### Database Migrations (Shared)
Database schema migrations are shared infrastructure.

**Location**: `infrastructure/database/migrations/`

```
infrastructure/
└── database/
    └── migrations/           # All database migrations
        ├── 000001_init.up.sql
        ├── 000001_init.down.sql
        ├── 000002_create_users_table.up.sql
        └── 000002_create_users_table.down.sql
```

**Note**: Database connection setup (`postgres.go`) is in `internal/platform/database/`, not here.

**Why migrations are shared?**
- ✅ Database is shared infrastructure
- ✅ Easier to manage all migrations in one place
- ✅ Single migration tool/runner

### SQLC Configuration (Single File)

**One `sqlc.yaml` at project root** generates code for all domains.

**Location**: `sqlc.yaml` (project root)

```yaml
# sqlc.yaml (at project root)
version: "2"
sql:
  # Users domain
  - engine: "postgresql"
    queries: "internal/domains/users/infrastructure/persistence/queries"
    schema: "infrastructure/database/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/domains/users/infrastructure/persistence/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
  
  # Products domain (example)
  - engine: "postgresql"
    queries: "internal/domains/products/infrastructure/persistence/queries"
    schema: "infrastructure/database/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/domains/products/infrastructure/persistence/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
```

**Generate code for all domains**:
```bash
# From project root
sqlc generate
```

**Why this approach?**
- ✅ Single command generates code for all domains
- ✅ Simpler for beginners (one config file)
- ✅ Easier to manage and maintain
- ✅ Each domain still has its own queries and generated code
- ✅ Clear ownership - queries live with domains

### PGX Connection Pool
- Single connection pool in `internal/platform/database/postgres.go`
- Injected into all domain repositories

## API Documentation (OpenAPI + Scalar)

### OpenAPI Specification

**Location**: `docs/openapi/openapi.yaml`

**Purpose**: API contract documentation using OpenAPI 3.0 specification.

```yaml
# docs/openapi/openapi.yaml
openapi: 3.0.0
info:
  title: Go DDD Clean Starter API
  version: 1.0.0
  description: RESTful API built with DDD and Clean Architecture
  contact:
    name: API Support
    email: support@example.com

servers:
  - url: http://localhost:3000/api/v1
    description: Development server
  - url: https://api.yourdomain.com/api/v1
    description: Production server

paths:
  /users:
    post:
      summary: Create a new user
      tags:
        - Users
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserResponse'
        '400':
          description: Invalid request
        '409':
          description: Email already exists

  /users/{id}:
    get:
      summary: Get user by ID
      tags:
        - Users
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: User found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserResponse'
        '404':
          description: User not found

components:
  schemas:
    CreateUserRequest:
      type: object
      required:
        - email
        - name
        - password
      properties:
        email:
          type: string
          format: email
          example: user@example.com
        name:
          type: string
          minLength: 2
          maxLength: 100
          example: John Doe
        password:
          type: string
          minLength: 8
          example: SecurePass123!
    
    UserResponse:
      type: object
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
          format: email
        name:
          type: string
        is_active:
          type: boolean
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time
```

### Scalar Integration

**Scalar** provides a beautiful, interactive API documentation UI powered by your OpenAPI spec.

#### Setup Scalar Endpoint

```go
// internal/platform/docs/handler.go
package docs

import (
    "github.com/gofiber/fiber/v2"
)

func RegisterDocsRoutes(app *fiber.App) {
    // Serve OpenAPI spec
    app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
        return c.SendFile("./docs/openapi/openapi.yaml")
    })
    
    // Serve Scalar UI
    app.Get("/docs", func(c *fiber.Ctx) error {
        html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
    <script 
        id="api-reference" 
        data-url="/openapi.yaml"
        data-configuration='{"theme":"purple"}'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
        c.Set("Content-Type", "text/html")
        return c.SendString(html)
    })
}
```

#### Register in main.go

```go
// cmd/api/main.go
package main

import (
    "github.com/gofiber/fiber/v2"
    "yourproject/internal/platform/docs"
    // ... other imports
)

func main() {
    app := fiber.New()
    
    // Register API documentation (development only)
    if config.Env == "development" {
        docs.RegisterDocsRoutes(app)
    }
    
    // ... rest of setup
    
    app.Listen(":3000")
}
```

#### Access Documentation

- **Scalar UI**: `http://localhost:3000/docs`
- **OpenAPI Spec**: `http://localhost:3000/openapi.yaml`

### Why OpenAPI + Scalar?

**Benefits**:
- ✅ **Interactive API Explorer**: Test endpoints directly from the browser
- ✅ **Type-Safe Contracts**: Define request/response schemas upfront
- ✅ **Client Generation**: Generate TypeScript/Python/etc. clients from spec
- ✅ **API-First Development**: Design API before implementation
- ✅ **Better than Swagger UI**: Faster, more modern, better UX
- ✅ **Documentation as Code**: Version controlled with your codebase

**Workflow**:
1. Define endpoint in `openapi.yaml`
2. Implement handler following the spec
3. Test via Scalar UI
4. Frontend/mobile teams use the spec to generate clients

### OpenAPI Best Practices

1. **Keep spec updated** - Update `openapi.yaml` when adding/changing endpoints
2. **Use components** - Define reusable schemas in `components/schemas`
3. **Add examples** - Include example requests/responses
4. **Document errors** - Specify all possible error responses
5. **Version your API** - Use `/api/v1`, `/api/v2` for breaking changes
6. **Validate spec** - Use `openapi-generator validate` or similar tools

### Alternative: Generate from Code

If you prefer code-first approach, you can use:
- **swaggo/swag** - Generate OpenAPI from Go comments
- **ogen** - Generate Go code from OpenAPI spec

However, for a **docs-first** starter, manually maintaining `openapi.yaml` is recommended as it encourages API design before implementation.

## Platform Layer (Technical Infrastructure)

**Location**: `internal/platform/`

**Purpose**: Technical infrastructure with **zero business logic**. Contains code that is technically necessary for the entire app.

```
internal/platform/
├── database/
│   └── postgres.go          # Database connection pool
├── logger/
│   └── logger.go            # Structured logging
├── middleware/
│   ├── auth.go              # Authentication middleware
│   ├── cors.go              # CORS middleware
│   ├── logger.go            # Request logging middleware
│   └── recovery.go          # Panic recovery
├── config/
│   └── config.go            # Application configuration
├── errors/
│   └── errors.go            # Common error types
├── docs/
│   └── handler.go           # API documentation handler (Scalar)
└── utils/
    └── validator.go         # Validation utilities
```

**What belongs in platform**:
- ✅ Database connection setup
- ✅ Logger configuration
- ✅ HTTP middleware
- ✅ Configuration loading
- ✅ Common utilities (validation, etc.)
- ✅ Error types for technical failures

**What does NOT belong in platform**:
- ❌ Business logic
- ❌ Domain models
- ❌ Use cases
- ❌ Domain-specific errors

## Simplified Project Structure

```
go-ddd-clean-starter/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point (Fiber setup)
│
├── internal/
│   ├── domains/                       # Domain silos
│   │   ├── users/                     # User domain
│   │   │   ├── domain/
│   │   │   │   ├── user.go           # Entity + value objects
│   │   │   │   ├── repository.go     # Repository interface
│   │   │   │   └── errors.go         # Domain errors
│   │   │   ├── application/
│   │   │   │   ├── service.go        # Use cases
│   │   │   │   └── dto.go            # DTOs
│   │   │   ├── infrastructure/
│   │   │   │   ├── persistence/
│   │   │   │   │   ├── queries/      # SQLC queries for users
│   │   │   │   │   │   └── users.sql
│   │   │   │   │   └── sqlc/         # Generated code for users
│   │   │   │   └── repository.go     # Repository implementation
│   │   │   └── handler/
│   │   │       ├── handler.go        # HTTP handlers
│   │   │       ├── request.go        # Request models
│   │   │       └── response.go       # Response models
│   │   │
│   │   └── products/                  # Product domain (example)
│   │       ├── domain/
│   │       ├── application/
│   │       ├── infrastructure/
│   │       │   └── persistence/
│   │       │       ├── queries/
│   │       │       └── sqlc/
│   │       └── handler/
│   │
│   └── platform/                     # Technical infrastructure (no business logic)
│       ├── database/
│       │   └── postgres.go           # DB connection pool
│       ├── logger/
│       │   └── logger.go
│       ├── middleware/
│       │   ├── auth.go
│       │   ├── cors.go
│       │   ├── logger.go
│       │   └── recovery.go
│       ├── config/
│       │   └── config.go
│       ├── errors/
│       │   └── errors.go
│       ├── docs/
│       │   └── handler.go            # API docs handler (Scalar)
│       └── utils/
│           └── validator.go
│
├── infrastructure/                   # Shared database migrations only
│   └── database/
│       └── migrations/               # All migrations
│           ├── 000001_create_users_table.up.sql
│           └── 000001_create_users_table.down.sql
│
├── docs/
│   ├── ARCHITECTURE.md
│   └── openapi/
│       └── openapi.yaml               # OpenAPI 3.0 specification
│
├── sqlc.yaml                          # SQLC config (generates for all domains)
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── README.md
```

**Clean and organized!**
- ✅ Each domain owns its queries and generated code
- ✅ Database migrations in shared infrastructure
- ✅ Platform layer for technical infrastructure only
- ✅ Clear separation of concerns
- ✅ Easy to understand and extend

## Key Architectural Principles

### 1. Domain Isolation
- Each domain is a **separate module** with its own layers
- Domains do NOT import each other
- Communication through events or APIs

### 2. Clean Architecture Per Domain
- **Domain layer**: Pure business logic (no dependencies)
- **Application layer**: Use cases (depends on domain only)
- **Infrastructure layer**: External concerns (implements interfaces)
- **Presentation layer**: HTTP handling (depends on application)

### 3. Dependency Rule
- Dependencies point **inward** (toward domain)
- Domain layer has **zero external dependencies**
- Infrastructure implements domain interfaces

### 4. Fiber is Routing Only
- Fiber **only** handles HTTP routing
- Fiber **does NOT** contain business logic
- Fiber **does NOT** call domain layer directly
- Handlers call application services

### 5. Database Infrastructure
- **Migrations in one place** (`infrastructure/database/migrations/`) - shared
- **Single SQLC config** at project root - generates code for all domains
- **Queries live with domain** (`internal/domains/users/infrastructure/persistence/queries/`)
- **Generated code stays in domain** (`internal/domains/users/infrastructure/persistence/sqlc/`)

### 6. Platform Layer (Technical Only)
- Contains **zero business logic**
- Database connection, logger, middleware, config, utils
- Shared across all domains
- Purely technical infrastructure

## Benefits of This Architecture

### Domain-Centric Benefits
1. **Domain Isolation**: Each domain can evolve independently
2. **Team Autonomy**: Teams can own entire domains
3. **Clear Boundaries**: No accidental coupling between domains
4. **Microservice Ready**: Easy to extract domains into services

### Clean Architecture Benefits
5. **Testability**: Domain logic testable without infrastructure
6. **Flexibility**: Easy to swap databases, frameworks
7. **Maintainability**: Clear separation of concerns
8. **Type Safety**: SQLC provides compile-time SQL validation

### Fiber Benefits
9. **Performance**: Fast HTTP routing
10. **Simplicity**: Fiber only does routing, nothing more
11. **No Framework Lock-in**: Business logic independent of Fiber

## Implementation Guide

### Step 1: Setup Shared Infrastructure
```bash
# Create database migrations folder
mkdir -p infrastructure/database/migrations

# Create platform layer
mkdir -p internal/platform/{database,logger,middleware,config,errors,docs,utils}

# Create OpenAPI documentation
mkdir -p docs/openapi
```

### Step 2: Create Domain Structure
```bash
mkdir -p internal/domains/users/{domain,application,handler}
mkdir -p internal/domains/users/infrastructure/persistence/{queries,sqlc}
```

### Step 3: Create SQLC Config (One Time)
```bash
# Create sqlc.yaml at project root (only need to do this once)
cat > sqlc.yaml << 'EOF'
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/domains/users/infrastructure/persistence/queries"
    schema: "infrastructure/database/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/domains/users/infrastructure/persistence/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
EOF
```

### Step 4: Domain Layer (Pure Business Logic)
- Define entity in `user.go`
- Define repository interface
- Define domain errors

### Step 5: Application Layer (Use Cases)
- Create service with use cases
- Define DTOs

### Step 6: Infrastructure Layer
- Write SQLC queries in `infrastructure/persistence/queries/users.sql`
- Run `sqlc generate` from project root to generate code
- Implement repository using generated SQLC code
- **Important**: Map SQLC models to domain entities (don't leak SQLC types)

### Step 7: Handler Layer
- Create HTTP handlers
- Define request/response models

### Step 8: Setup API Documentation
- Create `docs/openapi/openapi.yaml` with API specification
- Create Scalar handler in `internal/platform/docs/handler.go`
- Register docs routes in `main.go` (development only)

### Step 9: Wire in main.go
- Initialize database connection pool
- Initialize platform components (logger, config)
- Create domain services with dependency injection
- Register domain routes with Fiber
- Register API documentation routes (if development)
- Start server

## Example: Adding a New Domain

```bash
# 1. Add database migration
echo "CREATE TABLE products (...)" > infrastructure/database/migrations/000002_create_products_table.up.sql

# 2. Create domain structure
mkdir -p internal/domains/products/{domain,application,handler}
mkdir -p internal/domains/products/infrastructure/persistence/{queries,sqlc}

# 3. Add SQLC queries for products domain
echo "-- name: CreateProduct :one\nINSERT INTO products..." > internal/domains/products/infrastructure/persistence/queries/products.sql

# 4. Update sqlc.yaml at project root (add products domain)
cat >> sqlc.yaml << 'EOF'
  # Products domain
  - engine: "postgresql"
    queries: "internal/domains/products/infrastructure/persistence/queries"
    schema: "infrastructure/database/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/domains/products/infrastructure/persistence/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
EOF

# 5. Generate SQLC code for all domains
sqlc generate

# 6. Implement layers (domain -> application -> infrastructure -> handler)
# 7. Register routes in main.go
```

## Database Schema Template

### Migration (Shared Infrastructure)

```sql
-- infrastructure/database/migrations/000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### SQLC Queries (Per Domain)

```sql
-- internal/domains/users/infrastructure/persistence/queries/users.sql
-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND is_active = true;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND is_active = true;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE users SET is_active = false WHERE id = $1;
```

## Summary

This **simplified** architecture provides:
- ✅ **Domain-centric structure** (each domain is a silo)
- ✅ **Clean architecture inside each domain** (4 simple layers)
- ✅ **Fiber for routing only** (no business logic)
- ✅ **Shared database migrations** (SQLC per domain)
- ✅ **Minimal boilerplate** (perfect for a starter template)
- ✅ **Clear boundaries** (no cross-domain dependencies)
- ✅ **Easy to understand** (simple folder structure)
- ✅ **Easy to extend** (add new domains easily)

**Key Differences from Complex DDD**:
- Database migrations are **shared infrastructure**
- SQLC queries are **per domain** (domain owns its data access)
- Platform layer for technical infrastructure (no business logic)
- Simplified file structure (no excessive nesting)
- Fewer files per domain (merged related concerns)
- Perfect balance of **simplicity** and **clean architecture**

## Best Practices & Important Considerations

### 1. Context Propagation (Critical)

**Always pass `context.Context` as the first argument** in your application services and repository methods.

```go
// ✅ CORRECT
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserResponseDTO, error) {
    // ctx is passed through to repository
    user, err := s.repo.Save(ctx, domainUser)
    return result, err
}

// ✅ Repository also accepts context
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
    // pgx will use this context
    return r.queries.CreateUser(ctx, sqlc.CreateUserParams{...})
}
```

**Why this matters**:
- ✅ If user cancels HTTP request, database query is automatically cancelled
- ✅ Timeouts propagate through the entire call chain
- ✅ Saves database CPU cycles on cancelled requests
- ✅ Fiber provides `c.Context()` which is the request context

**In Handlers**:
```go
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // Pass Fiber's context to service
    result, err := h.service.CreateUser(c.Context(), dto)
    // ...
}
```

### 2. SQLC Model Mapping (Critical)

**Never expose SQLC-generated types outside the infrastructure layer.**

```go
// ❌ WRONG - Leaking SQLC types
func (r *UserRepository) FindByID(ctx context.Context, id string) (*sqlc.User, error) {
    return r.queries.GetUserByID(ctx, id) // Returns SQLC type directly
}

// ✅ CORRECT - Map to domain entity
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    sqlcUser, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Map SQLC model to domain entity
    return r.toDomain(sqlcUser), nil
}

func (r *UserRepository) toDomain(sqlcUser *sqlc.User) *domain.User {
    email, _ := domain.NewEmail(sqlcUser.Email)
    return &domain.User{
        ID:        sqlcUser.ID,
        Email:     email,
        Name:      sqlcUser.Name,
        IsActive:  sqlcUser.IsActive,
        CreatedAt: sqlcUser.CreatedAt,
        UpdatedAt: sqlcUser.UpdatedAt,
    }
}

func (r *UserRepository) toSQLC(user *domain.User) sqlc.CreateUserParams {
    return sqlc.CreateUserParams{
        Email:        user.Email.Value(),
        Name:         user.Name,
        PasswordHash: user.PasswordHash,
    }
}
```

**Why this matters**:
- ✅ If you change a column name in SQL, only the mapping code breaks (isolated)
- ✅ Domain layer stays pure and independent of database schema
- ✅ You can change database without affecting business logic
- ✅ Acts as an anti-corruption layer

### 3. Transaction Handling (Unit of Work)

**Transactions should be managed in the Application Layer**, not the infrastructure layer.

**Why?** The application service knows which operations need to be atomic.

#### Add Transaction Manager to Platform

```go
// internal/platform/database/transaction.go
package database

import (
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}

type txManager struct {
    pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) TxManager {
    return &txManager{pool: pool}
}

func (tm *txManager) WithTransaction(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
    tx, err := tm.pool.Begin(ctx)
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback(ctx)
            panic(p)
        }
    }()
    
    if err := fn(ctx, tx); err != nil {
        tx.Rollback(ctx)
        return err
    }
    
    return tx.Commit(ctx)
}
```

#### Use in Application Service

```go
// internal/domains/orders/application/service.go
type OrderService struct {
    orderRepo   domain.OrderRepository
    balanceRepo domain.BalanceRepository
    txManager   database.TxManager
}

func (s *OrderService) CreateOrder(ctx context.Context, dto CreateOrderDTO) error {
    // Start transaction
    return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
        // 1. Debit balance
        if err := s.balanceRepo.DebitWithTx(ctx, tx, dto.UserID, dto.Amount); err != nil {
            return err
        }
        
        // 2. Create order
        order := domain.NewOrder(dto.UserID, dto.ProductID, dto.Amount)
        if err := s.orderRepo.SaveWithTx(ctx, tx, order); err != nil {
            return err
        }
        
        // Both succeed or both rollback
        return nil
    })
}
```

#### Repository with Transaction Support

```go
// internal/domains/orders/infrastructure/repository.go
type OrderRepository struct {
    pool    *pgxpool.Pool
    queries *sqlc.Queries
}

// Regular method (no transaction)
func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
    params := r.toSQLC(order)
    _, err := r.queries.CreateOrder(ctx, params)
    return err
}

// Transaction-aware method
func (r *OrderRepository) SaveWithTx(ctx context.Context, tx pgx.Tx, order *domain.Order) error {
    // Create queries instance with transaction
    txQueries := r.queries.WithTx(tx)
    params := r.toSQLC(order)
    _, err := txQueries.CreateOrder(ctx, params)
    return err
}
```

**Key Points**:
- ✅ Application service decides what needs to be atomic
- ✅ Transaction manager lives in platform layer (technical concern)
- ✅ Repositories provide both regular and transaction-aware methods
- ✅ SQLC's `WithTx()` method makes this easy

### 4. Error Handling Strategy

**Map errors at layer boundaries**:

```go
// Infrastructure layer
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrUserNotFound  // Map to domain error
        }
        return nil, fmt.Errorf("database error: %w", err)
    }
    return r.toDomain(sqlcUser), nil
}

// Handler layer
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    user, err := h.service.GetUser(c.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrUserNotFound) {
            return c.Status(404).JSON(fiber.Map{"error": "user not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
    }
    return c.JSON(user)
}
```

### 5. SQLC Configuration Best Practice

**Single `sqlc.yaml` at project root** - simple and clean:

```yaml
# sqlc.yaml (at project root)
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/domains/users/infrastructure/persistence/queries"
    schema: "infrastructure/database/migrations"
    gen:
      go:
        package: "sqlc"
        out: "internal/domains/users/infrastructure/persistence/sqlc"
        sql_package: "pgx/v5"
```

**Generate code**:
```bash
# From project root
sqlc generate
```

**Benefits**:
- ✅ No complex relative paths
- ✅ One command generates all domains
- ✅ Easy to add new domains (just add new entry)
- ✅ Simpler for beginners
