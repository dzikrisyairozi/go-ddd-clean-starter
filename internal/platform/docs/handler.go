package docs

import (
	"github.com/gofiber/fiber/v2"
)

/*
RegisterRoutes registers API documentation endpoints with the Fiber app.
It provides two routes:
  - GET /openapi.yaml - Serves the OpenAPI 3.0 specification file
  - GET /docs - Serves an interactive Scalar UI for browsing and testing the API

The Scalar UI is a modern, fast alternative to Swagger UI with better UX.
These routes should typically only be enabled in development environments.
*/
func RegisterRoutes(app *fiber.App) {
	// Serve OpenAPI specification
	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/openapi/openapi.yaml")
	})

	// Serve Scalar API documentation UI
	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation - Go DDD Clean Starter</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <script 
        id="api-reference" 
        data-url="/openapi.yaml"
        data-configuration='{
            "theme": "purple",
            "layout": "modern",
            "defaultOpenAllTags": true
        }'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})
}
