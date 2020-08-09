package accounts

import "errors"

// Authenticate makes sure the credentials are valid in the DB, and then returns the account's ID.
// If the email doesn't exist, it returns an EmailNotExistsError.
// If the password is wrong, it returns an InvalidPasswordError.
func (s *PostgresStore) Authenticate(email, password string) (int64, error) {
	return -1, errors.New("not yet implemented")
}
