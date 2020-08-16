package accounts

import (
	"context"
	"errors"
)

// See the docs on interfaces in store.go
func (s *PostgresStore) SetForgottenPassword(ctx context.Context, id int64, password, resetToken string) error {
	return errors.New("not yet implemented")
}

// See the docs on interfaces in store.go
func (s *PostgresStore) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	return errors.New("not yet implemented")
}
