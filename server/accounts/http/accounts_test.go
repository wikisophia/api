package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestAccountAcceptsEmails(t *testing.T) {
	app := acceptancetest.NewApp(t, nil)
	rr := doSaveAccount(app, `{"email":"some-email@soph.wiki"}`)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	welcomeEmails := app.Emailer.Welcomes
	require.Len(t, welcomeEmails, 1)
	require.Equal(t, "some-email@soph.wiki", welcomeEmails[0].Email)
	require.Len(t, welcomeEmails[0].ResetToken, 20) // Tokens need to be long enough for security
}

func TestAccountSendsResetEmails(t *testing.T) {
	app := acceptancetest.NewApp(t, nil)
	assert.Equal(t, http.StatusNoContent, doSaveAccount(app, `{"email":"some-email@soph.wiki"}`).Code)
	assert.Equal(t, http.StatusNoContent, doSaveAccount(app, `{"email":"some-email@soph.wiki"}`).Code)
	require.Len(t, app.Emailer.Welcomes, 1)
	require.Len(t, app.Emailer.PasswordResets, 1)
	require.Equal(t, app.Emailer.Welcomes[0].ID, app.Emailer.PasswordResets[0].ID)
	require.Equal(t, app.Emailer.Welcomes[0].Email, app.Emailer.PasswordResets[0].Email)
}

func TestAccountRejectsBadRequestBodies(t *testing.T) {
	acceptancetest.AssertBadRequest(t, "POST", "/accounts", "not json")
	acceptancetest.AssertBadRequest(t, "POST", "/accounts", "{}")
	acceptancetest.AssertBadRequest(t, "POST", "/accounts", `{"email":null}`)
	acceptancetest.AssertBadRequest(t, "POST", "/accounts", `{"email":5}`)
	acceptancetest.AssertBadRequest(t, "POST", "/accounts", `{"email":true}`)
	acceptancetest.AssertBadRequest(t, "POST", "/accounts", `{"email":3.4}`)
}

func doSaveAccount(a *acceptancetest.App, body string) *httptest.ResponseRecorder {
	return a.Do(httptest.NewRequest("POST", "/accounts", strings.NewReader(body)))
}
