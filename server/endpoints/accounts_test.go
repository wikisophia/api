package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/endpoints"
)

func TestAccountAcceptsEmails(t *testing.T) {
	app := newAppForTests(t, nil)
	rr := doSaveAccount(app.server, `{"email":"some-email@soph.wiki"}`)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	welcomeEmails := app.emailer.welcomes
	require.Len(t, welcomeEmails, 1)
	require.Equal(t, "some-email@soph.wiki", welcomeEmails[0].Email)
	require.Len(t, welcomeEmails[0].ResetToken, 20) // Tokens need to be long enough for security
}

func TestAccountSendsResetEmails(t *testing.T) {
	app := newAppForTests(t, nil)
	assert.Equal(t, http.StatusNoContent, doSaveAccount(app.server, `{"email":"some-email@soph.wiki"}`).Code)
	assert.Equal(t, http.StatusNoContent, doSaveAccount(app.server, `{"email":"some-email@soph.wiki"}`).Code)
	require.Len(t, app.emailer.welcomes, 1)
	require.Len(t, app.emailer.passwordResets, 1)
	require.Equal(t, app.emailer.welcomes[0].ID, app.emailer.passwordResets[0].ID)
	require.Equal(t, app.emailer.welcomes[0].Email, app.emailer.passwordResets[0].Email)
}

func TestAccountRejectsBadRequestBodies(t *testing.T) {
	assertBadRequest(t, doSaveAccount(newAppForTests(t, nil).server, "not json"))
	assertBadRequest(t, doSaveAccount(newAppForTests(t, nil).server, "{}"))
	assertBadRequest(t, doSaveAccount(newAppForTests(t, nil).server, `{"email":null}`))
	assertBadRequest(t, doSaveAccount(newAppForTests(t, nil).server, `{"email":5}`))
	assertBadRequest(t, doSaveAccount(newAppForTests(t, nil).server, `{"email":true}`))
	assertBadRequest(t, doSaveAccount(newAppForTests(t, nil).server, `{"email":3.4}`))
}

func doSaveAccount(s *endpoints.Server, body string) *httptest.ResponseRecorder {
	return doRequest(s, httptest.NewRequest("POST", "/accounts", strings.NewReader(body)))
}

func assertBadRequest(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
