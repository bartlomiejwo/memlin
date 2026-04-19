package db

import (
	"backend/internal/config"
	sqlcdb "backend/internal/db/sqlc"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*sqlcdb.Queries
	pool *pgxpool.Pool
}

// Connect establishes a connection to the database and returns a Store
func Connect(cfg *config.Config) (*Store, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	// Create a new store with the pgx pool
	return &Store{
		Queries: sqlcdb.New(pool),
		pool:    pool,
	}, nil
}

// Close closes the database connection
func (s *Store) Close() {
	s.pool.Close()
}

// GetPool returns the pgxpool.Pool instance
func (s *Store) GetPool() *pgxpool.Pool {
	return s.pool
}

// WithTx executes a function within a transaction
func (s *Store) WithTx(ctx context.Context, fn func(*sqlcdb.Queries) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := sqlcdb.New(tx)
	err = fn(q)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
