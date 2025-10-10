package handler

import (
	"errors"
	"net/http"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/application"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/domain"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

/*
UserHandler handles HTTP requests for user endpoints.
This is the outermost layer that deals with HTTP concerns only.
It delegates business logic to the application service.
*/
type UserHandler struct {
	userService *application.UserService
	logger      *logger.Logger
}

/*
NewUserHandler creates a new UserHandler instance.
Requires a UserService for business logic and a Logger for request logging.
*/
func NewUserHandler(userService *application.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      log,
	}
}

/*
CreateUser handles POST /users - Create a new user.
Request body: CreateUserRequest
Response: 201 Created with UserResponse
Errors: 400 Bad Request, 409 Conflict (email exists), 500 Internal Server Error
*/
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Convert to DTO
	dto := application.CreateUserDTO{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	// Call service
	user, err := h.userService.CreateUser(c.Context(), dto)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(http.StatusCreated).JSON(toUserResponse(user))
}

/*
GetUser handles GET /users/:id - Get a user by ID.
Path parameter: id (UUID)
Response: 200 OK with UserResponse
Errors: 400 Bad Request (invalid ID), 404 Not Found, 500 Internal Server Error
*/
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Parse ID from path
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
	}

	// Call service
	user, err := h.userService.GetUser(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(http.StatusOK).JSON(toUserResponse(user))
}

/*
UpdateUser handles PUT /users/:id - Update a user.
Path parameter: id (UUID)
Request body: UpdateUserRequest
Response: 200 OK with UserResponse
Errors: 400 Bad Request, 404 Not Found, 409 Conflict, 500 Internal Server Error
*/
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Parse ID from path
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
	}

	// Parse request body
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Convert to DTO
	dto := application.UpdateUserDTO{
		Email: req.Email,
		Name:  req.Name,
	}

	// Call service
	user, err := h.userService.UpdateUser(c.Context(), id, dto)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(http.StatusOK).JSON(toUserResponse(user))
}

/*
DeleteUser handles DELETE /users/:id - Delete a user (soft delete).
Path parameter: id (UUID)
Response: 204 No Content
Errors: 400 Bad Request, 404 Not Found, 500 Internal Server Error
*/
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Parse ID from path
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
	}

	// Call service
	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}

	// Return no content
	return c.SendStatus(http.StatusNoContent)
}

/*
ListUsers handles GET /users - List users with pagination.
Query parameters: limit (default 10, max 100), offset (default 0)
Response: 200 OK with UserListResponse
Errors: 400 Bad Request, 500 Internal Server Error
*/
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Parse query parameters
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Validate parameters
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Call service
	users, err := h.userService.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(http.StatusOK).JSON(toUserListResponse(users))
}

/*
ChangePassword handles POST /users/:id/password - Change user password.
Path parameter: id (UUID)
Request body: ChangePasswordRequest
Response: 200 OK with success message
Errors: 400 Bad Request, 401 Unauthorized (wrong old password), 404 Not Found, 500 Internal Server Error
*/
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	// Parse ID from path
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
	}

	// Parse request body
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Convert to DTO
	dto := application.ChangePasswordDTO{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	// Call service
	if err := h.userService.ChangePassword(c.Context(), id, dto); err != nil {
		return h.handleError(c, err)
	}

	// Return success
	return c.Status(http.StatusOK).JSON(SuccessResponse{
		Message: "Password changed successfully",
	})
}

/*
handleError maps domain/application errors to appropriate HTTP status codes.
This is the error translation layer between business logic and HTTP.
*/
func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	// Log the error
	h.logger.Error("Handler error", "error", err.Error(), "path", c.Path())

	// Map domain errors to HTTP status codes
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return c.Status(http.StatusNotFound).JSON(ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})

	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return c.Status(http.StatusConflict).JSON(ErrorResponse{
			Error:   "conflict",
			Message: "Email already exists",
		})

	case errors.Is(err, domain.ErrInvalidEmail):
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_email",
			Message: "Invalid email format",
		})

	case errors.Is(err, domain.ErrInvalidPassword):
		return c.Status(http.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "invalid_password",
			Message: "Invalid password",
		})

	case errors.Is(err, domain.ErrUserInactive):
		return c.Status(http.StatusForbidden).JSON(ErrorResponse{
			Error:   "user_inactive",
			Message: "User account is inactive",
		})

	default:
		// Generic internal server error
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "An internal error occurred",
		})
	}
}
