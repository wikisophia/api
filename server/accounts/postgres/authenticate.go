package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

const authenticateQuery = `
SELECT id, password_hash
FROM accounts
WHERE email = $1;
`

// See the docs on interfaces in store.go
func (s *PostgresStore) Authenticate(ctx context.Context, email, password string) (int64, error) {
	row := s.pool.QueryRow(ctx, authenticateQuery, email)
	var id int64
	var hashedPassword string
	if err := row.Scan(&id, &hashedPassword); err == pgx.ErrNoRows {
		return -1, errors.New("no account found with that email and password")
	} else if err != nil {
		return -1, err
	}

	match, err := s.hasher.Matches(password, hashedPassword)
	if err != nil {
		return -1, errors.New("error matching password against the database")
	}
	if !match {
		return -1, errors.New("no account found with that email and password")
	}
	return id, nil
}
