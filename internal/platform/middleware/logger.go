package middleware

import (
	"time"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

/*
RequestLogger returns a Fiber middleware that logs all HTTP requests.
For each request, it:
  - Generates a unique request ID (stored in context as "requestID")
  - Records the request start time
  - Processes the request
  - Logs request details including method, path, status, duration, and client IP
  - Uses different log levels based on response status:
  - ERROR (red) for 5xx server errors
  - WARN (yellow) for 4xx client errors
  - INFO (green) for 2xx/3xx successful responses

The request ID can be retrieved in handlers via c.Locals("requestID") for correlation.
*/
func RequestLogger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Locals("requestID", requestID)

		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		// Choose log level based on status code
		if status >= 500 {
			log.Error("HTTP Request",
				"request_id", requestID,
				"method", method,
				"path", path,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", ip,
			)
		} else if status >= 400 {
			log.Warn("HTTP Request",
				"request_id", requestID,
				"method", method,
				"path", path,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", ip,
			)
		} else {
			log.Info("HTTP Request",
				"request_id", requestID,
				"method", method,
				"path", path,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"ip", ip,
			)
		}

		return err
	}
}
