package accounts

// EmailExistsError will be returned if callers try to create a new account
// with an email that already exists in the system.
type EmailExistsError struct {
	Email string
}

func (e EmailExistsError) Error() string {
	return e.Email + " already has an account"
}

// CorruptedPasswordError will be returned if the hashed password is corrupted for some reason.
// This really shouldn't happen, unless the code is buggy or DB hardware has issues.
type CorruptedPasswordError struct {
	Email string
}

func (e CorruptedPasswordError) Error() string {
	return "the password for " + e.Email + " has been corrupted in the storage backend"
}

// EmailNotExistsError will be returned if the user tried to operate on an email address
// which doesn't have an associated Account in the system.
type EmailNotExistsError struct {
	Email string
}

func (e EmailNotExistsError) Error() string {
	return e.Email + " does not have an account"
}

// InvalidPasswordError will be returned if the user authenticated with the wrong password.
type InvalidPasswordError struct{}

func (i InvalidPasswordError) Error() string {
	return "invalid password"
}

// ProhibitedPasswordError will be returned if the user tries to set a password which
// we don't allow.
type ProhibitedPasswordError struct{}

func (e ProhibitedPasswordError) Error() string {
	return "the password is unacceptable"
}

// InvalidResetTokenError will be returned if the user sent an unrecognized
// or expired password reset token when trying to reset their password.
type InvalidResetTokenError struct{}

func (err InvalidResetTokenError) Error() string {
	return "unrecognized verification token"
}
