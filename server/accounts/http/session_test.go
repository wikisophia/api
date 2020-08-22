package http_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestSessionsErrorCodes(t *testing.T) {
	acceptancetest.AssertMethodNotAllowed(t, http.MethodConnect, "/sessions")
	acceptancetest.AssertMethodNotAllowed(t, http.MethodDelete, "/sessions")
	acceptancetest.AssertMethodNotAllowed(t, http.MethodGet, "/sessions")
	acceptancetest.AssertMethodNotAllowed(t, http.MethodPatch, "/sessions")
	acceptancetest.AssertMethodNotAllowed(t, http.MethodPut, "/sessions")
	acceptancetest.AssertMethodNotAllowed(t, http.MethodTrace, "/sessions")

	acceptancetest.AssertBadRequest(t, "POST", "/sessions", "")
	acceptancetest.AssertBadRequest(t, "POST", "/sessions", "not json")
	acceptancetest.AssertBadRequest(t, "POST", "/sessions", "5")
	acceptancetest.AssertBadRequest(t, "POST", "/sessions", "{}")
	acceptancetest.AssertBadRequest(t, "POST", "/sessions", `{"email":"something@soph.wiki"}`)
	acceptancetest.AssertBadRequest(t, "POST", "/sessions", `{"password":"password"}`)
}

func TestUnknownUserForbidden(t *testing.T) {
	assert.Equal(t, http.StatusForbidden, newApp(t, nil).Authenticate("something@soph.wiki", "some-password").Code)
}

func TestValidCredentialsAccepted(t *testing.T) {
	a := newApp(t, nil)
	a.SaveAccountWithPasswordSuccessfully("some-email@soph.wiki", "some-password")
	a.AuthenticateSuccessfully("some-email@soph.wiki", "some-password")
}

func TestInvalidCredentialsForbidden(t *testing.T) {
	a := newApp(t, nil)
	a.SaveAccountWithPasswordSuccessfully("some-email@soph.wiki", "some-password")
	assert.Equal(t, http.StatusForbidden, a.Authenticate("some-email@soph.wiki", "wrong-password").Code)
}
