package database

import (
	"context"
	"fmt"
	"time"

	"github.com/dzikrisyairozi/go-ddd-clean-starter/internal/platform/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
NewPool creates and configures a new PostgreSQL connection pool using pgxpool.
It parses the database URL from configuration, sets pool parameters (max/min connections,
lifetime, idle time), and verifies connectivity with a ping.
Returns an error if configuration parsing, pool creation, or database ping fails.
The pool should be closed when the application shuts down using Close().
*/
func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// Build connection config
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = cfg.Database.MaxConns
	poolConfig.MinConns = cfg.Database.MinConns
	poolConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

/*
Close gracefully shuts down the database connection pool.
It waits for all active connections to finish their work before closing.
Should be called during application shutdown, typically in a defer statement.
*/
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
