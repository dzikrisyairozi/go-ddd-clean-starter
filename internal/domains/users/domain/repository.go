package domain

import (
	"context"

	"github.com/google/uuid"
)

/*
UserRepository defines the contract for user persistence operations.
This is a domain interface (port) that will be implemented by the infrastructure layer.
The domain layer defines what it needs, and the infrastructure layer provides it.

This follows the Dependency Inversion Principle:
  - High-level domain logic does not depend on low-level infrastructure details
  - Both depend on this abstraction (interface)

All methods accept context.Context as the first parameter for:
  - Request cancellation propagation
  - Timeout handling
  - Request-scoped values (like request ID)
*/
type UserRepository interface {
	/*
		Save persists a new user to the repository.
		Returns an error if:
		  - The email already exists (ErrEmailAlreadyExists)
		  - Database connection fails
		  - Validation fails

		The user ID should be generated before calling Save.
	*/
	Save(ctx context.Context, user *User) error

	/*
		FindByID retrieves a user by their unique identifier.
		Returns:
		  - The user if found
		  - ErrUserNotFound if no user exists with the given ID
		  - Other errors for database failures

		Only returns active users by default (is_active = true).
	*/
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	/*
		FindByEmail retrieves a user by their email address.
		Returns:
		  - The user if found
		  - ErrUserNotFound if no user exists with the given email
		  - Other errors for database failures

		Email lookup is case-insensitive.
		Only returns active users by default (is_active = true).
	*/
	FindByEmail(ctx context.Context, email Email) (*User, error)

	/*
		Update modifies an existing user in the repository.
		Returns an error if:
		  - The user does not exist (ErrUserNotFound)
		  - The new email conflicts with another user (ErrEmailAlreadyExists)
		  - Database connection fails
		  - Validation fails

		The UpdatedAt timestamp should be set before calling Update.
	*/
	Update(ctx context.Context, user *User) error

	/*
		Delete removes a user from the repository (soft delete).
		This sets is_active = false rather than physically deleting the record.
		Returns an error if:
		  - The user does not exist (ErrUserNotFound)
		  - Database connection fails

		Soft delete allows for:
		  - Data retention for auditing
		  - Potential account recovery
		  - Maintaining referential integrity
	*/
	Delete(ctx context.Context, id uuid.UUID) error

	/*
		List retrieves all active users with pagination support.
		Parameters:
		  - limit: Maximum number of users to return
		  - offset: Number of users to skip (for pagination)

		Returns:
		  - Slice of users (may be empty)
		  - Error if database operation fails

		Only returns active users (is_active = true).
		Results are ordered by created_at DESC (newest first).
	*/
	List(ctx context.Context, limit, offset int) ([]*User, error)

	/*
		Count returns the total number of active users.
		Useful for pagination calculations.
		Returns:
		  - Count of active users
		  - Error if database operation fails
	*/
	Count(ctx context.Context) (int64, error)
}
