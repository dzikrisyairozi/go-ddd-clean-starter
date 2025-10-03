package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

/*
User represents a user entity in the domain.
This is the core business object with identity and lifecycle.
Contains business rules and behavior related to users.
*/
type User struct {
	ID           uuid.UUID
	Email        Email
	Name         string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

/*
Email is a value object representing a validated email address.
Value objects are immutable and defined by their attributes.
Email validation is enforced at creation time.
*/
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

/*
NewEmail creates a new Email value object with validation.
Returns an error if the email format is invalid.
Email is converted to lowercase for consistency.
Example: NewEmail("User@Example.com") returns Email{value: "user@example.com"}
*/
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, ErrInvalidEmail
	}

	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: email}, nil
}

/*
Value returns the string representation of the email.
This is the only way to access the email value, maintaining encapsulation.
*/
func (e Email) Value() string {
	return e.value
}

/*
String implements the Stringer interface for Email.
Allows Email to be printed directly.
*/
func (e Email) String() string {
	return e.value
}

/*
NewUser creates a new User entity with the provided details.
This is a factory function that enforces business rules:
  - Email must be valid
  - Name must not be empty
  - Password hash must not be empty
  - New users are active by default
  - Timestamps are set to current time

Returns an error if any validation fails.
Example:

	email, _ := NewEmail("user@example.com")
	user, err := NewUser(email, "John Doe", "hashed_password")
*/
func NewUser(email Email, name, passwordHash string) (*User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	now := time.Now()

	return &User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

/*
UpdateProfile updates the user's name and email.
This method enforces business rules:
  - Name must not be empty
  - Email must be valid
  - UpdatedAt timestamp is automatically updated

Returns an error if validation fails.
*/
func (u *User) UpdateProfile(name string, email Email) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now()

	return nil
}

/*
Deactivate marks the user as inactive (soft delete).
This implements the soft delete pattern - the user record remains in the database
but is marked as inactive and won't appear in normal queries.
UpdatedAt timestamp is automatically updated.
*/
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

/*
Activate marks the user as active.
Used to restore a previously deactivated user.
UpdatedAt timestamp is automatically updated.
*/
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

/*
ChangePassword updates the user's password hash.
The password should already be hashed before calling this method.
UpdatedAt timestamp is automatically updated.
Returns an error if the password hash is empty.
*/
func (u *User) ChangePassword(passwordHash string) error {
	if passwordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()

	return nil
}

/*
Validate checks if the user entity is in a valid state.
This is useful before persisting the user to the database.
Returns an error describing what is invalid.
*/
func (u *User) Validate() error {
	if u.ID == uuid.Nil {
		return errors.New("user ID cannot be nil")
	}

	if u.Email.Value() == "" {
		return errors.New("email cannot be empty")
	}

	if u.Name == "" {
		return errors.New("name cannot be empty")
	}

	if u.PasswordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	return nil
}

/*
String returns a string representation of the user.
Useful for logging and debugging.
Does not include sensitive information like password hash.
*/
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Email: %s, Name: %s, IsActive: %t}",
		u.ID, u.Email, u.Name, u.IsActive)
}
