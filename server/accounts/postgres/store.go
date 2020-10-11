package postgres

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgresStore returns a Store which can manage accounts.
// The returned Store.Close() function will *not* close the pool, since we did not open it.
func NewPostgresStore(pool *pgxpool.Pool, hasher Hasher) *PostgresStore {
	if pool == nil {
		log.Fatal("A connection pool is required to make an accounts.PostgresStore.")
	}
	return &PostgresStore{
		pool:   pool,
		hasher: hasher,
	}
}

type Hasher interface {
	// Hash a value to a string which encodes the algorithm + salt as well.
	Hash(value string) (string, error)
	// Check if the given value matches a hash.
	Matches(value string, hash string) (bool, error)
}

// PostgresStore saves account info in Postgres.
type PostgresStore struct {
	pool   *pgxpool.Pool
	hasher Hasher
}
