package application

import (
	"time"

	"github.com/google/uuid"
)

/*
DTOs (Data Transfer Objects) for the application layer.
These objects transfer data between layers without exposing domain entities directly.
They provide a stable API contract independent of domain model changes.
*/

// CreateUserDTO represents the input data for creating a new user
type CreateUserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"` // Plain text password (will be hashed)
}

// UpdateUserDTO represents the input data for updating a user
type UpdateUserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ChangePasswordDTO represents the input data for changing a user's password
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserResponseDTO represents the output data for a user
// This is what gets returned to clients (handlers, APIs)
type UserResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponseDTO represents a paginated list of users
type UserListResponseDTO struct {
	Users   []UserResponseDTO `json:"users"`
	Total   int64             `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"has_more"`
}
