package handler

import (
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/application"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/logger"
	"github.com/gofiber/fiber/v2"
)

/*
RegisterRoutes registers all user-related routes with the Fiber app.
This follows the principle: "Fiber for routing ONLY".
All business logic is in the application service.

Routes:

	POST   /users           - Create a new user
	GET    /users           - List users (paginated)
	GET    /users/:id       - Get a user by ID
	PUT    /users/:id       - Update a user
	DELETE /users/:id       - Delete a user (soft delete)
	POST   /users/:id/password - Change user password
*/
func RegisterRoutes(app *fiber.App, userService *application.UserService, log *logger.Logger) {
	// Create handler
	handler := NewUserHandler(userService, log)

	// User routes
	users := app.Group("/users")

	users.Post("/", handler.CreateUser)                 // Create user
	users.Get("/", handler.ListUsers)                   // List users
	users.Get("/:id", handler.GetUser)                  // Get user by ID
	users.Put("/:id", handler.UpdateUser)               // Update user
	users.Delete("/:id", handler.DeleteUser)            // Delete user
	users.Post("/:id/password", handler.ChangePassword) // Change password
}
