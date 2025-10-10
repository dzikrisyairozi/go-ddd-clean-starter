package handler

/*
Request models for user endpoints.
These structs define the expected JSON structure for incoming HTTP requests.
They are separate from DTOs to allow for different validation rules at the HTTP layer.
*/

// CreateUserRequest represents the request body for creating a new user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=1,max=100"`
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ListUsersQuery represents query parameters for listing users
type ListUsersQuery struct {
	Limit  int `query:"limit" validate:"min=1,max=100"`
	Offset int `query:"offset" validate:"min=0"`
}
