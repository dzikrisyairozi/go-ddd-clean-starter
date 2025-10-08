-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    name,
    password_hash,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = true
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = true
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    name = $3,
    password_hash = $4,
    is_active = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET
    is_active = false,
    updated_at = $2
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE is_active = true;
