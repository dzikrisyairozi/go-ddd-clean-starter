package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/domain"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/infrastructure/persistence/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
UserRepository implements the domain.UserRepository interface using SQLC.
This is the infrastructure layer implementation that handles actual database operations.
It depends on SQLC-generated code for type-safe database queries.
*/
type UserRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

/*
NewUserRepository creates a new UserRepository instance.
Requires a pgxpool.Pool for database connectivity.
The SQLC Queries instance is created from the pool.
*/
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

/*
Save persists a new user to the database.
Maps the domain User entity to SQLC parameters and executes the insert query.
Returns ErrEmailAlreadyExists if a user with the same email already exists.
*/
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	params := sqlc.CreateUserParams{
		ID:           uuidToPgtype(user.ID),
		Email:        user.Email.Value(),
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		IsActive:     user.IsActive,
		CreatedAt:    timeToPgtype(user.CreatedAt),
		UpdatedAt:    timeToPgtype(user.UpdatedAt),
	}

	_, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		// Check for unique constraint violation (email already exists)
		if isUniqueViolation(err) {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

/*
FindByID retrieves a user by their unique identifier.
Returns ErrUserNotFound if no active user exists with the given ID.
Maps the SQLC User model to a domain User entity.
*/
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	sqlcUser, err := r.queries.GetUserByID(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.toDomainUser(sqlcUser)
}

/*
FindByEmail retrieves a user by their email address.
Returns ErrUserNotFound if no active user exists with the given email.
Email lookup is case-insensitive (handled by database).
*/
func (r *UserRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email.Value())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.toDomainUser(sqlcUser)
}

/*
Update modifies an existing user in the database.
Maps the domain User entity to SQLC parameters and executes the update query.
Returns ErrUserNotFound if the user doesn't exist.
Returns ErrEmailAlreadyExists if the new email conflicts with another user.
*/
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	params := sqlc.UpdateUserParams{
		ID:           uuidToPgtype(user.ID),
		Email:        user.Email.Value(),
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		IsActive:     user.IsActive,
		UpdatedAt:    timeToPgtype(user.UpdatedAt),
	}

	_, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		if isUniqueViolation(err) {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

/*
Delete soft deletes a user by setting is_active = false.
The user record remains in the database for auditing purposes.
Returns ErrUserNotFound if the user doesn't exist.
*/
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	params := sqlc.DeleteUserParams{
		ID:        uuidToPgtype(id),
		UpdatedAt: timeToPgtype(time.Now()),
	}

	err := r.queries.DeleteUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

/*
List retrieves a paginated list of active users.
Results are ordered by created_at DESC (newest first).
Maps SQLC User models to domain User entities.
*/
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	params := sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	sqlcUsers, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		user, err := r.toDomainUser(sqlcUser)
		if err != nil {
			return nil, fmt.Errorf("failed to map user at index %d: %w", i, err)
		}
		users[i] = user
	}

	return users, nil
}

/*
Count returns the total number of active users.
Useful for pagination calculations.
*/
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

/*
toDomainUser maps a SQLC User model to a domain User entity.
This is the anti-corruption layer that prevents database models from leaking into the domain.
*/
func (r *UserRepository) toDomainUser(sqlcUser sqlc.User) (*domain.User, error) {
	email, err := domain.NewEmail(sqlcUser.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return &domain.User{
		ID:           pgtypeToUUID(sqlcUser.ID),
		Email:        email,
		Name:         sqlcUser.Name,
		PasswordHash: sqlcUser.PasswordHash,
		IsActive:     sqlcUser.IsActive,
		CreatedAt:    pgtypeToTime(sqlcUser.CreatedAt),
		UpdatedAt:    pgtypeToTime(sqlcUser.UpdatedAt),
	}, nil
}

/*
isUniqueViolation checks if the error is a PostgreSQL unique constraint violation.
This is used to detect email conflicts.
*/
func isUniqueViolation(err error) bool {
	// PostgreSQL error code 23505 is unique_violation
	// This is a simplified check - in production you might want to use pgconn.PgError
	return err != nil && (err.Error() == "ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)" ||
		contains(err.Error(), "unique constraint") ||
		contains(err.Error(), "duplicate key"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Type conversion helpers

/*
uuidToPgtype converts uuid.UUID to pgtype.UUID.
This is needed because SQLC uses pgtype.UUID for PostgreSQL UUID columns.
*/
func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

/*
pgtypeToUUID converts pgtype.UUID to uuid.UUID.
This is needed to convert SQLC types back to domain types.
*/
func pgtypeToUUID(pgUUID pgtype.UUID) uuid.UUID {
	return pgUUID.Bytes
}

/*
timeToPgtype converts time.Time to pgtype.Timestamp.
This is needed because SQLC uses pgtype.Timestamp for PostgreSQL TIMESTAMP columns.
*/
func timeToPgtype(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

/*
pgtypeToTime converts pgtype.Timestamp to time.Time.
This is needed to convert SQLC types back to domain types.
*/
func pgtypeToTime(pgTime pgtype.Timestamp) time.Time {
	return pgTime.Time
}
