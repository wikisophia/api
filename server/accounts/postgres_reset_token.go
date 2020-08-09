package accounts

import "errors"

// NewResetTokenWithAccount gets a password reset token for the given email.
// If the account doesn't exist yet, it will be created first.
func (s *PostgresStore) NewResetTokenWithAccount(email string) (string, error) {
	return "", errors.New("not yet implemented")
}
