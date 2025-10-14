package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/application"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/handler"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/domains/users/infrastructure/persistence"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/config"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/database"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/docs"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/logger"
	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize logger
	log := logger.New("info")
	log.Info("Starting Go DDD Clean Starter API...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err.Error())
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration", "error", err.Error())
	}

	log.Info("Configuration loaded successfully",
		"env", cfg.App.Environment,
		"port", cfg.App.Port,
		"db_host", cfg.Database.Host)

	// Initialize database connection pool
	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err.Error())
	}
	defer pool.Close()

	log.Info("Database connection established successfully")

	// Initialize dependencies (Dependency Injection)
	// Infrastructure layer
	userRepo := persistence.NewUserRepository(pool)

	// Application layer
	userService := application.NewUserService(userRepo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Go DDD Clean Starter API",
		ServerHeader: "Fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   "internal_error",
				"message": err.Error(),
			})
		},
	})

	// Register middleware
	app.Use(middleware.Recovery(log))
	app.Use(middleware.RequestLogger(log))
	app.Use(middleware.CORS())

	// Register API documentation routes
	docs.RegisterDocsRoutes(app)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Register domain routes
	handler.RegisterRoutes(app, userService, log)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Info("Shutting down server gracefully...")

		// Shutdown with timeout
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Error("Server forced to shutdown", "error", err.Error())
		}

		// Close database connection
		pool.Close()
		log.Info("Server stopped")
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Info("Server starting", "address", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server", "error", err.Error())
	}
}
