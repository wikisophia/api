package accounts

import (
	"context"
	"errors"
)

// See the docs on interfaces in store.go
func (s *PostgresStore) Authenticate(ctx context.Context, email, password string) (int64, error) {
	return -1, errors.New("not yet implemented")
}
