package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
UserService implements the application use cases for the users domain.
It orchestrates domain objects and coordinates business workflows.
This layer depends only on the domain layer (not on infrastructure).

The service:
  - Validates input data
  - Coordinates domain entities
  - Defines transaction boundaries
  - Maps between DTOs and domain entities
  - Handles password hashing (application concern, not domain)
*/
type UserService struct {
	userRepo domain.UserRepository
}

/*
NewUserService creates a new UserService instance.
Requires a UserRepository implementation (provided by infrastructure layer).
This follows dependency injection pattern.
*/
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

/*
CreateUser creates a new user account.
This use case:
 1. Validates input data
 2. Checks if email already exists
 3. Hashes the password
 4. Creates domain entity
 5. Persists to repository

Returns the created user or an error if:
  - Email is invalid
  - Email already exists
  - Password is too weak
  - Repository operation fails
*/
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserResponseDTO, error) {
	// Validate input
	if err := s.validateCreateUserInput(dto); err != nil {
		return nil, err
	}

	// Create email value object
	email, err := domain.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := s.hashPassword(dto.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create domain entity
	user, err := domain.NewUser(email, dto.Name, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Persist to repository
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Map to response DTO
	return s.toUserResponseDTO(user), nil
}

/*
GetUser retrieves a user by ID.
Returns the user or ErrUserNotFound if not found.
*/
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*UserResponseDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toUserResponseDTO(user), nil
}

/*
GetUserByEmail retrieves a user by email address.
Returns the user or ErrUserNotFound if not found.
*/
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserResponseDTO, error) {
	emailVO, err := domain.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, emailVO)
	if err != nil {
		return nil, err
	}

	return s.toUserResponseDTO(user), nil
}

/*
UpdateUser updates a user's profile information.
This use case:
 1. Retrieves the existing user
 2. Validates new data
 3. Checks if new email conflicts with another user
 4. Updates the domain entity
 5. Persists changes

Returns the updated user or an error.
*/
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*UserResponseDTO, error) {
	// Retrieve existing user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate input
	if dto.Name == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Create new email value object
	newEmail, err := domain.NewEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	// Check if new email conflicts with another user
	if newEmail.Value() != user.Email.Value() {
		existingUser, err := s.userRepo.FindByEmail(ctx, newEmail)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, domain.ErrEmailAlreadyExists
		}
	}

	// Update domain entity
	if err := user.UpdateProfile(dto.Name, newEmail); err != nil {
		return nil, err
	}

	// Persist changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponseDTO(user), nil
}

/*
DeleteUser soft deletes a user account.
The user record remains in the database but is marked as inactive.
*/
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

/*
ActivateUser activates a previously deactivated user account.
*/
func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) (*UserResponseDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Activate()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	return s.toUserResponseDTO(user), nil
}

/*
DeactivateUser deactivates a user account.
*/
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) (*UserResponseDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Deactivate()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to deactivate user: %w", err)
	}

	return s.toUserResponseDTO(user), nil
}

/*
ChangePassword changes a user's password.
Verifies the old password before setting the new one.
*/
func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, dto ChangePasswordDTO) error {
	// Retrieve user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := s.verifyPassword(user.PasswordHash, dto.OldPassword); err != nil {
		return domain.ErrInvalidPassword
	}

	// Validate new password
	if err := s.validatePassword(dto.NewPassword); err != nil {
		return err
	}

	// Hash new password
	newPasswordHash, err := s.hashPassword(dto.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := user.ChangePassword(newPasswordHash); err != nil {
		return err
	}

	// Persist changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

/*
ListUsers retrieves a paginated list of active users.
*/
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) (*UserListResponseDTO, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	// Get users
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Map to response DTOs
	userDTOs := make([]UserResponseDTO, len(users))
	for i, user := range users {
		userDTOs[i] = *s.toUserResponseDTO(user)
	}

	return &UserListResponseDTO{
		Users:   userDTOs,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+limit) < total,
	}, nil
}

// Helper methods

func (s *UserService) validateCreateUserInput(dto CreateUserDTO) error {
	if dto.Email == "" {
		return errors.New("email is required")
	}
	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Password == "" {
		return errors.New("password is required")
	}
	return s.validatePassword(dto.Password)
}

func (s *UserService) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// Add more password validation rules as needed
	return nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *UserService) verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *UserService) toUserResponseDTO(user *domain.User) *UserResponseDTO {
	return &UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email.Value(),
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
