package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

/*
CORS returns a Fiber middleware that handles Cross-Origin Resource Sharing (CORS).
Configured with permissive defaults suitable for development:
  - Allows all origins (*)
  - Allows common HTTP methods (GET, POST, PUT, DELETE, PATCH, OPTIONS)
  - Allows standard headers (Origin, Content-Type, Accept, Authorization)
  - Cache preflight requests for 24 hours

For production, use CORSWithConfig() to specify allowed origins explicitly.
*/
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 24 hours
	})
}

/*
CORSWithConfig returns a CORS middleware with custom configuration.
Allows fine-grained control over CORS behavior including:
  - Specific allowed origins (e.g., "https://yourdomain.com")
  - Allowed HTTP methods
  - Allowed and exposed headers
  - Credentials support
  - Preflight cache duration

Use this in production to restrict access to trusted origins only.
*/
func CORSWithConfig(config cors.Config) fiber.Handler {
	return cors.New(config)
}
