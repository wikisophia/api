package postgres

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgresStore returns a Store which is used to save and load Arguments.
// The db should point to a Postgres database.
// The returned Store.Close() function will *not* close this connection, since we did not open it.
func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	if pool == nil {
		log.Fatal("A connection pool is required to make an arguments.PostgresStore.")
	}
	return &PostgresStore{
		pool: pool,
	}
}

// PostgresStore expects that {projectRoot}/postgres/scripts/create.sql
// has already been run on your database so that the expected schema exists.
type PostgresStore struct {
	pool *pgxpool.Pool
}
