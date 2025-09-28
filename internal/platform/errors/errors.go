package errors

import (
	"errors"
	"fmt"
)

// Common error types for platform-level errors
var (
	// Database errors
	ErrDatabaseConnection  = errors.New("database connection error")
	ErrDatabaseQuery       = errors.New("database query error")
	ErrDatabaseTransaction = errors.New("database transaction error")

	// Validation errors
	ErrValidation   = errors.New("validation error")
	ErrInvalidInput = errors.New("invalid input")

	// Configuration errors
	ErrConfiguration = errors.New("configuration error")
	ErrMissingConfig = errors.New("missing configuration")

	// General errors
	ErrInternal     = errors.New("internal server error")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// AppError represents an application error with additional context
type AppError struct {
	Err     error
	Message string
	Code    string
	Fields  map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

/*
Unwrap returns the underlying error wrapped by AppError.
This allows errors.Is() and errors.As() to work correctly with wrapped errors.
*/
func (e *AppError) Unwrap() error {
	return e.Err
}

/*
NewAppError creates a new application error with an underlying error and custom message.
The message provides additional context about where or why the error occurred.
Example: NewAppError(err, "failed to create user")
*/
func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Fields:  make(map[string]interface{}),
	}
}

/*
WithCode adds an error code to the AppError for categorization.
Error codes can be used by clients to handle specific error types programmatically.
Returns the AppError for method chaining.
Example: NewAppError(err, "validation failed").WithCode("VALIDATION_ERROR")
*/
func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

/*
WithField adds a single key-value pair to the error's contextual information.
Useful for adding debugging information like IDs, timestamps, or other relevant data.
Returns the AppError for method chaining.
Example: err.WithField("user_id", 123)
*/
func (e *AppError) WithField(key string, value interface{}) *AppError {
	e.Fields[key] = value
	return e
}

/*
WithFields adds multiple key-value pairs to the error's contextual information.
Accepts a map of field names to values for bulk addition of context.
Returns the AppError for method chaining.
Example: err.WithFields(map[string]interface{}{"user_id": 123, "action": "create"})
*/
func (e *AppError) WithFields(fields map[string]interface{}) *AppError {
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

// ValidationError represents a validation error with field-level details
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

/*
NewValidationError creates a new validation error for a specific field.
Used to indicate that a particular field failed validation.
Example: NewValidationError("email", "must be a valid email address")
*/
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

// Error implements the error interface
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %d errors", len(e.Errors))
}

/*
Add appends a new validation error to the collection.
Used to accumulate multiple field validation errors before returning them.
Example: errs.Add("email", "is required")
*/
func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

/*
HasErrors returns true if the collection contains any validation errors.
Useful for checking if validation passed before proceeding with an operation.
*/
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

/*
NewValidationErrors creates a new empty collection for accumulating validation errors.
Use this at the start of validation logic, then call Add() for each validation failure.
Example: errs := NewValidationErrors(); if email == "" { errs.Add("email", "is required") }
*/
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}
