package domain

import "errors"

/*
Domain-specific errors for the users domain.
These errors represent business rule violations and domain-level failures.
They should be used by the domain layer and mapped to appropriate HTTP status codes
in the handler layer.
*/

var (
	// ErrUserNotFound indicates that a user with the given identifier does not exist
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailAlreadyExists indicates that a user with the given email already exists
	// This enforces the business rule that emails must be unique
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidEmail indicates that the provided email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidPassword indicates that the provided password does not meet requirements
	ErrInvalidPassword = errors.New("invalid password")

	// ErrUserInactive indicates that the user account is deactivated
	// Operations on inactive users may be restricted
	ErrUserInactive = errors.New("user is inactive")

	// ErrUnauthorized indicates that the user is not authorized to perform the action
	ErrUnauthorized = errors.New("unauthorized")
)
