package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager manages database transactions
type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}

// txManager implements TxManager
type txManager struct {
	pool *pgxpool.Pool
}

/*
NewTxManager creates a new transaction manager instance.
The transaction manager provides a convenient way to execute database operations
within a transaction with automatic commit/rollback handling.
Requires a pgxpool.Pool instance for database connectivity.
*/
func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &txManager{pool: pool}
}

/*
WithTransaction executes a function within a database transaction.
It automatically handles transaction lifecycle:
  - Begins a new transaction
  - Executes the provided function with the transaction context
  - Commits if the function succeeds
  - Rolls back if the function returns an error
  - Rolls back and re-panics if a panic occurs

This ensures atomic operations - either all database changes succeed or none do.
Example usage:

	txManager.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
	    // Perform multiple database operations
	    return nil
	})
*/
func (tm *txManager) WithTransaction(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	// Begin transaction
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of panic
	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = tx.Rollback(ctx)
			panic(p) // Re-throw panic after rollback
		}
	}()

	// Execute function
	if err := fn(ctx, tx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
