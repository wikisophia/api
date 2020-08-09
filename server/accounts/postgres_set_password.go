package accounts

import (
	"context"
	"errors"
)

// SetPassword associates the password with this email, if the resetToken is valid.
// If the email doesn't exist, it returns an EmailNotExistsError.
// If the resetToken is wrong (expired or never returned by ResetPassword(email)),
//   it returns an InvalidPasswordError.
func (s *PostgresStore) SetPassword(ctx context.Context, email, password, resetToken string) (int64, error) {
	return -1, errors.New("not yet implemented")
}
