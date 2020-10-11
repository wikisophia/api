package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

const authenticateQuery = `
SELECT id
FROM accounts
WHERE email = $1
	AND password_hash = $2;
`

// See the docs on interfaces in store.go
func (s *PostgresStore) Authenticate(ctx context.Context, email, password string) (int64, error) {
	row := s.pool.QueryRow(ctx, authenticateQuery, email, password)
	var id int64
	if err := row.Scan(&id); err == pgx.ErrNoRows {
		return -1, errors.New("no account found with that email and password")
	} else if err != nil {
		return -1, err
	}
	return id, nil
}
