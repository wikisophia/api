package accounts

import "strconv"

// EmailExistsError will be returned if callers try to create a new account
// with an email that already exists in the system.
type EmailExistsError struct {
	Email string
}

func (e *EmailExistsError) Error() string {
	return e.Email + " already has an account"
}

// CorruptedPasswordError will be returned if the hashed password is corrupted for some reason.
// This really shouldn't happen, and probably indicates a bug.
type CorruptedPasswordError struct {
	Email string
}

func (e *CorruptedPasswordError) Error() string {
	return "the password for " + e.Email + " has been corrupted in the storage backend"
}

// EmailNotExistsError will be returned if callers try to authenticate with an
// email that doesn't exist in teh system.
type EmailNotExistsError struct {
	Email string
}

func (e *EmailNotExistsError) Error() string {
	return e.Email + " does not have an account"
}

// InvalidPasswordError will be returned if the password doesn't match.
type InvalidPasswordError struct{}

func (i InvalidPasswordError) Error() string {
	return "invalid password"
}

// AccountNotFoundError will be returned if someone asks the store
// for a user which doesn't exist.
type AccountNotFoundError struct {
	ID int64
}

func (e *AccountNotFoundError) Error() string {
	return "user " + strconv.FormatInt(e.ID, 10) + " does not exist"
}

// The ExpiredVerificationTokenError signals that the verification link sent to
// the user's email expired by the time they clicked on it.
type ExpiredVerificationTokenError struct{}

func (err ExpiredVerificationTokenError) Error() string {
	return "expired verification token"
}

// The InvalidVerificationTokenError signals that the user sent an unrecognized token.
type InvalidVerificationTokenError struct{}

func (err InvalidVerificationTokenError) Error() string {
	return "unrecognized verification token"
}
