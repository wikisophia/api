package postgres

import (
	"context"
	"errors"

	"github.com/wikisophia/api/server/accounts"
)

// See the docs on interfaces in store.go
func (s *PostgresStore) NewResetToken(ctx context.Context, email string) (accounts.Account, bool, error) {
	return accounts.Account{}, false, errors.New("not yet implemented")
}
