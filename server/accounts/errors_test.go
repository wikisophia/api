package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/accounts"
)

func TestErrorMessages(t *testing.T) {
	assert.EqualError(t,
		accounts.EmailExistsError{"some-mail@soph.wiki"},
		"some-mail@soph.wiki already has an account")
	assert.EqualError(t,
		accounts.CorruptedPasswordError{"some-mail@soph.wiki"},
		"the password for some-mail@soph.wiki has been corrupted in the storage backend")
	assert.EqualError(t,
		accounts.AccountNotExistsError{"some-mail@soph.wiki"},
		"some-mail@soph.wiki does not have an account")
	assert.EqualError(t, accounts.InvalidPasswordError{}, "invalid password")
	assert.EqualError(t, accounts.ProhibitedPasswordError{}, "the password is unacceptable")
	assert.EqualError(t, accounts.InvalidResetTokenError{}, "unrecognized verification token")
}
