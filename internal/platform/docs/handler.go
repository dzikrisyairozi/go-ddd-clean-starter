package docs

import (
	"github.com/gofiber/fiber/v2"
)

/*
RegisterDocsRoutes registers API documentation routes.
This provides both the OpenAPI spec and Scalar UI for interactive API documentation.

Routes:

	GET /openapi.yaml - OpenAPI specification file
	GET /docs         - Scalar interactive documentation UI
*/
func RegisterDocsRoutes(app *fiber.App) {
	// Serve OpenAPI specification
	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/openapi/openapi.yaml")
	})

	// Serve Scalar UI
	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation - Go DDD Clean Starter</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
    <script 
        id="api-reference" 
        data-url="/openapi.yaml"
        data-configuration='{"theme":"purple","layout":"modern"}'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})
}
