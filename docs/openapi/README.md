# OpenAPI Modular Structure

This directory contains a modular OpenAPI 3.0 specification organized by domain.

## Structure

```
docs/openapi/
├── openapi.yaml                    # Main entry point
├── components/                     # Shared/common components
│   ├── parameters.yaml            # Shared parameters
│   ├── responses.yaml             # Shared responses
│   └── schemas/                   # Shared schemas
│       ├── error.yaml
│       └── success.yaml
└── domains/                        # Domain-specific specs
    └── users/                      # Users domain
        ├── paths/                  # Endpoint definitions
        │   ├── users.yaml         # POST /users, GET /users
        │   ├── user-by-id.yaml    # GET/PUT/DELETE /users/:id
        │   └── change-password.yaml # POST /users/:id/password
        └── schemas/                # Request/response schemas
            ├── create-user-request.yaml
            ├── update-user-request.yaml
            ├── change-password-request.yaml
            ├── user-response.yaml
            └── user-list-response.yaml
```

## Benefits

### 1. **Domain-Driven Organization**
- Each domain has its own directory
- Easy to find and maintain domain-specific specs
- Follows the same structure as your code

### 2. **Scalability**
- Add new domains without bloating the main file
- Each domain team can work independently
- Hundreds of endpoints stay organized

### 3. **Reusability**
- Shared components in `/components`
- Domain-specific schemas in domain folders
- No duplication

### 4. **Maintainability**
- Small, focused files
- Easy to review changes
- Clear separation of concerns

## Adding a New Domain

When adding a new domain (e.g., "products"):

```
1. Create domain directory:
   docs/openapi/domains/products/

2. Add path files:
   docs/openapi/domains/products/paths/
   ├── products.yaml
   └── product-by-id.yaml

3. Add schema files:
   docs/openapi/domains/products/schemas/
   ├── create-product-request.yaml
   ├── product-response.yaml
   └── product-list-response.yaml

4. Reference in main openapi.yaml:
   paths:
     /products:
       $ref: './domains/products/paths/products.yaml'
```

## How It Works

OpenAPI 3.0 supports `$ref` to reference external files:

```yaml
# Main file references domain paths
paths:
  /users:
    $ref: './domains/users/paths/users.yaml'

# Path file references schemas
requestBody:
  content:
    application/json:
      schema:
        $ref: '../schemas/create-user-request.yaml'

# Schemas can reference other schemas
properties:
  users:
    type: array
    items:
      $ref: './user-response.yaml'
```

## Viewing the Docs

The modular structure is transparent to API consumers. When you visit:

```
http://localhost:3000/docs
```

Scalar UI automatically resolves all `$ref` and shows a unified, complete API documentation.

## Best Practices

1. **One endpoint per path file** (or group related methods like GET/PUT/DELETE)
2. **One schema per file**
3. **Shared components in `/components`**
4. **Domain-specific in `/domains/{domain}`**
5. **Use relative paths for `$ref`**

## Example: Adding a New Endpoint

To add `GET /users/:id/orders`:

```yaml
# 1. Create path file
# docs/openapi/domains/users/paths/user-orders.yaml
get:
  tags:
    - Users
  summary: Get user orders
  parameters:
    - $ref: '../../components/parameters.yaml#/UserId'
  responses:
    '200':
      content:
        application/json:
          schema:
            $ref: '../schemas/user-orders-response.yaml'

# 2. Create schema file
# docs/openapi/domains/users/schemas/user-orders-response.yaml
type: object
properties:
  orders:
    type: array
    items:
      type: object

# 3. Reference in main openapi.yaml
paths:
  /users/{id}/orders:
    $ref: './domains/users/paths/user-orders.yaml'
```

Done! The new endpoint appears in the docs automatically.
