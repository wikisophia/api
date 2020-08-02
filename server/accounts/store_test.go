package accounts_test

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api/server/accounts"
)

// StoreTests is a testing suite which makes sure that a Store obeys
// the interface contract
type StoreTests struct {
	suite.Suite
	StoreFactory func() Store
}

// Store has all the functions needed to manage accounts.
type Store interface {
	NewAccount(email string) (string, error)
	ResetPassword(email string) (string, error)
	SetPassword(email, password, resetToken string) (int64, error)
	Authenticate(email, password string) (int64, error)
}

// TestNewUserFlow makes sure users can make a new account, set their password, and log in.
func (suite *StoreTests) TestNewUserFlow() {
	store := suite.StoreFactory()
	token, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	id, err := store.SetPassword("email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)

	authed, err := store.Authenticate("email@soph.wiki", "password")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, authed)
}

// TestEmptyPasswordFailsBeforeInitialSet makes sure people can't log in with an
// empty password after they've created an account, but before setting it the first time.
func (suite *StoreTests) TestEmptyPasswordFailsBeforeInitialSet() {
	store := suite.StoreFactory()
	_, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.Authenticate("email@soph.wiki", "")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
}

// TestEmptyPasswordsRejected makes sure the user can't set their password to be empty.
func (suite *StoreTests) TestEmptyPasswordsRejected() {
	store := suite.StoreFactory()
	token, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword("email@soph.wiki", "", token)
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.ProhibitedPasswordError{}))
}

// TestInvalidTokensRejected makes sure people can't set their password with the wrong reset token.
func (suite *StoreTests) TestInvalidTokensRejected() {
	store := suite.StoreFactory()
	token, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword("email@soph.wiki", "password", token+"-wrong")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

// TestInvalidPasswordsRejected makes sure people can't log in with the wrong password.
func (suite *StoreTests) TestInvalidPasswordsRejected() {
	store := suite.StoreFactory()
	token, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword("email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)

	_, err = store.Authenticate("email@soph.wiki", "password2")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
}

func (suite *StoreTests) TestRegisterEmailTwiceReturnsError() {
	store := suite.StoreFactory()
	_, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.NewAccount("email@soph.wiki")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.EmailExistsError{}))
}

func (suite *StoreTests) TestTokenBecomesInvalidAfterUse() {
	store := suite.StoreFactory()
	token, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword("email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)
	_, err = store.SetPassword("email@soph.wiki", "password2", token)
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

func (suite *StoreTests) TestNewResetTokenInvalidatesOld() {
	store := suite.StoreFactory()
	token, err := store.NewAccount("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword("email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)
	firstResetToken, err := store.ResetPassword("email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.ResetPassword("email@soph.wiki")
	require.NoError(suite.T(), err)

	_, err = store.SetPassword("email@soph.wiki", "password2", firstResetToken)
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

func (suite *StoreTests) TestResetPasswordUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	_, err := store.ResetPassword("email@soph.wiki")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.EmailNotExistsError{}))
}

func (suite *StoreTests) TestSetPasswordUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	_, err := store.SetPassword("email@soph.wiki", "password", "token")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.EmailNotExistsError{}))
}

func (suite *StoreTests) TestAuthenticateUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	_, err := store.Authenticate("email@soph.wiki", "password")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.EmailNotExistsError{}))
}
