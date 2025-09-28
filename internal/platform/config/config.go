package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string
	Port        string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level string
}

/*
Load reads configuration from environment variables and .env file.
It loads application settings, database connection parameters, and logger configuration.
All values have sensible defaults if environment variables are not set.
Returns an error if required configuration values are missing or invalid.
*/
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "3000"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "go_ddd_starter"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxConns:        getEnvAsInt32("DB_MAX_CONNS", 25),
			MinConns:        getEnvAsInt32("DB_MIN_CONNS", 5),
			MaxConnLifetime: getEnvAsDuration("DB_MAX_CONN_LIFETIME", "1h"),
			MaxConnIdleTime: getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", "30m"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

/*
Validate checks if the configuration is valid by ensuring all required fields are present.
It verifies that critical database connection parameters (host, port, user, database name)
and application settings (port) are not empty.
Returns an error describing which required field is missing.
*/
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if c.App.Port == "" {
		return fmt.Errorf("application port is required")
	}
	return nil
}

/*
IsDevelopment returns true if the application is running in development mode.
This is determined by checking if APP_ENV is set to "development".
Useful for enabling development-only features like verbose logging or API documentation.
*/
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

/*
IsProduction returns true if the application is running in production mode.
This is determined by checking if APP_ENV is set to "production".
Useful for enabling production-only features like strict error handling or performance optimizations.
*/
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

/*
GetDatabaseURL constructs and returns a PostgreSQL connection string in the format:
postgres://user:password@host:port/database?sslmode=mode
This URL is used by pgx to establish database connections.
*/
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt32(key string, defaultValue int32) int32 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(value)
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		valueStr = defaultValue
	}
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		// Fallback to default if parsing fails
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}
