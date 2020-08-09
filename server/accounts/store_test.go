package accounts_test

import (
	"context"
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
	StoreFactory func() accounts.Store
}

// TestNewUserFlow makes sure users can make a new account, set their password, and log in.
func (suite *StoreTests) TestNewUserFlow() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	id, err := store.SetPassword(context.Background(), "email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)
	authed, err := store.Authenticate(context.Background(), "email@soph.wiki", "password")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, authed)
}

// TestEmptyPasswordFailsBeforeInitialSet makes sure people can't log in with an
// empty password after they've created an account, but before setting it the first time.
func (suite *StoreTests) TestEmptyPasswordFailsBeforeInitialSet() {
	store := suite.StoreFactory()
	_, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
}

// TestOldPasswordStillWorksAfterResetRequested makes sure the user can still use their old
// password after a reset has been requested.
func (suite *StoreTests) TestOldPasswordStillWorksAfterResetRequested() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "some-password", token)
	require.NoError(suite.T(), err)
	_, err = store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "some-password")
	require.NoError(suite.T(), err)
}

// TestResetPasswordInvalidatesOldOne makes sure the user's new password works, and
// the old password _doesn't_ work, after it's been reset.
func (suite *StoreTests) TestResetPasswordInvalidatesOldOne() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "some-password", token)
	require.NoError(suite.T(), err)
	token, err = store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "some-new-password", token)
	require.NoError(suite.T(), err)
	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "some-password")
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "some-new-password")
	require.NoError(suite.T(), err)
}

// TestSecondResetTokenInvalidatesFirst makes sure that the first reset token is
// invalidated after a second one is requested
func (suite *StoreTests) TestSecondResetTokenInvalidatesFirst() {
	store := suite.StoreFactory()
	token1, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	token2, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "some-password", token1)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "some-password", token2)
	require.NoError(suite.T(), err)
}

// TestEmptyPasswordsRejected makes sure the user can't set their password to be empty.
func (suite *StoreTests) TestEmptyPasswordsRejected() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "", token)
	require.True(suite.T(), errors.As(err, &accounts.ProhibitedPasswordError{}))
}

// TestInvalidTokensRejected makes sure people can't set their password with the wrong reset token.
func (suite *StoreTests) TestInvalidTokensRejected() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "password", "wrong-"+token)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

// TestInvalidPasswordsRejected makes sure people can't log in with the wrong password.
func (suite *StoreTests) TestInvalidPasswordsRejected() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)

	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "wrong-password")
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
}

func (suite *StoreTests) TestTokenBecomesInvalidAfterUse() {
	store := suite.StoreFactory()
	token, err := store.NewResetTokenWithAccount(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "password", token)
	require.NoError(suite.T(), err)
	_, err = store.SetPassword(context.Background(), "email@soph.wiki", "another-password", token)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

func (suite *StoreTests) TestSetPasswordUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	_, err := store.SetPassword(context.Background(), "email@soph.wiki", "password", "token")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.EmailNotExistsError{}))
}

func (suite *StoreTests) TestAuthenticateUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	_, err := store.Authenticate(context.Background(), "email@soph.wiki", "password")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.EmailNotExistsError{}))
}
