package postgres

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgresStore returns a Store which can manage accounts.
// The db should point to a Postgres database.
// The returned Store.Close() function will *not* close this connection, since we did not open it.
func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	if pool == nil {
		log.Fatal("A connection pool is required to make an accounts.PostgresStore.")
	}
	return &PostgresStore{
		pool: pool,
	}
}

// PostgresStore saves account info in Postgres.
type PostgresStore struct {
	pool *pgxpool.Pool
}
