package accounts

import (
	"context"
	"errors"
)

// See the docs on interfaces in store.go
func (s *PostgresStore) NewResetToken(ctx context.Context, email string) (Account, error) {
	return Account{}, errors.New("not yet implemented")
}
