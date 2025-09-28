package middleware

import (
	"fmt"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/logger"
	"github.com/gofiber/fiber/v2"
)

/*
Recovery returns a Fiber middleware that recovers from panics in request handlers.
When a panic occurs:
  - The panic is caught and logged with request details (method, path, request ID)
  - A 500 Internal Server Error response is sent to the client
  - The server continues running (prevents crash)

This middleware should be registered early in the middleware chain to catch
panics from all subsequent handlers. Critical for production stability.
*/
func Recovery(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic
				requestID, _ := c.Locals("requestID").(string)
				log.Error("Panic recovered",
					"request_id", requestID,
					"method", c.Method(),
					"path", c.Path(),
					"panic", fmt.Sprintf("%v", r),
				)

				// Return 500 Internal Server Error
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
			}
		}()

		return c.Next()
	}
}
