package handler

import (
	"time"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/application"
	"github.com/google/uuid"
)

/*
Response models for user endpoints.
These structs define the JSON structure for HTTP responses.
They map directly from application DTOs but can be customized for API needs.
*/

// UserResponse represents a single user in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users   []UserResponse `json:"users"`
	Total   int64          `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
	HasMore bool           `json:"has_more"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper functions to convert from application DTOs to response models

func toUserResponse(dto *application.UserResponseDTO) UserResponse {
	return UserResponse{
		ID:        dto.ID,
		Email:     dto.Email,
		Name:      dto.Name,
		IsActive:  dto.IsActive,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func toUserListResponse(dto *application.UserListResponseDTO) UserListResponse {
	users := make([]UserResponse, len(dto.Users))
	for i, user := range dto.Users {
		users[i] = toUserResponse(&user)
	}

	return UserListResponse{
		Users:   users,
		Total:   dto.Total,
		Limit:   dto.Limit,
		Offset:  dto.Offset,
		HasMore: dto.HasMore,
	}
}
