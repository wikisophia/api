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
	account, accountIsNew, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	require.True(suite.T(), accountIsNew)
	err = store.SetForgottenPassword(context.Background(), account.ID, "password", account.ResetToken)
	require.NoError(suite.T(), err)
	authed, err := store.Authenticate(context.Background(), "email@soph.wiki", "password")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.ID, authed)
}

// TestEmptyPasswordFailsBeforeInitialSet makes sure people can't log in with an
// empty password after they've created an account, but before setting it the first time.
func (suite *StoreTests) TestEmptyPasswordFailsBeforeInitialSet() {
	store := suite.StoreFactory()
	_, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
}

// TestOldPasswordStillWorksAfterResetRequested makes sure the user can still use their old
// password after a reset has been requested.
func (suite *StoreTests) TestOldPasswordStillWorksAfterResetRequested() {
	store := suite.StoreFactory()
	account, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "some-password", account.ResetToken)
	require.NoError(suite.T(), err)
	_, accountIsNew, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	require.False(suite.T(), accountIsNew)
	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "some-password")
	require.NoError(suite.T(), err)
}

// TestResetPasswordInvalidatesOldOne makes sure the user's new password works, and
// the old password _doesn't_ work, after it's been reset.
func (suite *StoreTests) TestResetPasswordInvalidatesOldOne() {
	store := suite.StoreFactory()
	account, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "some-password", account.ResetToken)
	require.NoError(suite.T(), err)
	account, _, err = store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "some-new-password", account.ResetToken)
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
	account1, account1IsNew, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	require.True(suite.T(), account1IsNew)
	account2, account2IsNew, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.Equal(suite.T(), account1.ID, account2.ID)
	require.Equal(suite.T(), account1.Email, account2.Email)
	require.False(suite.T(), account2IsNew)
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account1.ID, "some-password", account1.ResetToken)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
	err = store.SetForgottenPassword(context.Background(), account1.ID, "some-password", account2.ResetToken)
	require.NoError(suite.T(), err)
}

func (suite *StoreTests) TestDifferentAccountsAreIndependent() {
	store := suite.StoreFactory()
	firstAccount, firstIsNew, err := store.NewResetToken(context.Background(), "email1@soph.wiki")
	require.True(suite.T(), firstIsNew)
	require.NoError(suite.T(), err)
	secondAccount, secondIsNew, err := store.NewResetToken(context.Background(), "email2@soph.wiki")
	require.True(suite.T(), secondIsNew)
	require.NoError(suite.T(), err)

	err = store.SetForgottenPassword(context.Background(), firstAccount.ID, "some-password", secondAccount.ResetToken)
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
	err = store.SetForgottenPassword(context.Background(), secondAccount.ID, "some-password", firstAccount.ResetToken)
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

// TestEmptyPasswordsRejected makes sure the user can't set their password to be empty.
func (suite *StoreTests) TestEmptyPasswordsRejected() {
	store := suite.StoreFactory()
	account, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "", account.ResetToken)
	require.True(suite.T(), errors.As(err, &accounts.ProhibitedPasswordError{}))
}

// TestInvalidTokensRejected makes sure people can't set their password with the wrong reset token.
func (suite *StoreTests) TestInvalidTokensRejected() {
	store := suite.StoreFactory()
	account, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "password", "wrong-"+account.ResetToken)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

// TestInvalidPasswordsRejected makes sure people can't log in with the wrong password.
func (suite *StoreTests) TestInvalidPasswordsRejected() {
	store := suite.StoreFactory()
	account, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "password", account.ResetToken)
	require.NoError(suite.T(), err)

	_, err = store.Authenticate(context.Background(), "email@soph.wiki", "wrong-password")
	require.True(suite.T(), errors.As(err, &accounts.InvalidPasswordError{}))
}

func (suite *StoreTests) TestTokenBecomesInvalidAfterUse() {
	store := suite.StoreFactory()
	account, _, err := store.NewResetToken(context.Background(), "email@soph.wiki")
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "password", account.ResetToken)
	require.NoError(suite.T(), err)
	err = store.SetForgottenPassword(context.Background(), account.ID, "another-password", account.ResetToken)
	require.True(suite.T(), errors.As(err, &accounts.InvalidResetTokenError{}))
}

func (suite *StoreTests) TestSetPasswordUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	err := store.SetForgottenPassword(context.Background(), 1, "password", "token")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.AccountNotExistsError{}))
}

func (suite *StoreTests) TestAuthenticateUnknownEmailReturnsError() {
	store := suite.StoreFactory()
	_, err := store.Authenticate(context.Background(), "email@soph.wiki", "password")
	require.Error(suite.T(), err)
	require.True(suite.T(), errors.As(err, &accounts.AccountNotExistsError{}))
}
