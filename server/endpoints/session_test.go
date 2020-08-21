package endpoints_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/endpoints"
)

func TestSessionsErrorCodes(t *testing.T) {
	assertMethodNotAllowed(t, http.MethodConnect, "/sessions")
	assertMethodNotAllowed(t, http.MethodDelete, "/sessions")
	assertMethodNotAllowed(t, http.MethodGet, "/sessions")
	assertMethodNotAllowed(t, http.MethodPatch, "/sessions")
	assertMethodNotAllowed(t, http.MethodPut, "/sessions")
	assertMethodNotAllowed(t, http.MethodTrace, "/sessions")

	assertBadRequest(t, doAuthenticate(newAppForTests(t, nil).server, ""))
	assertBadRequest(t, doAuthenticate(newAppForTests(t, nil).server, "5"))
	assertBadRequest(t, doAuthenticate(newAppForTests(t, nil).server, "{}"))
	assertBadRequest(t, doAuthenticate(newAppForTests(t, nil).server, `{"email":"something@soph.wiki"}`))
	assertBadRequest(t, doAuthenticate(newAppForTests(t, nil).server, `{"password":"password"}`))

	assertBadRequest(t, doAuthenticate(newAppForTests(t, nil).server, `{"password":"password"}`))
}

func TestUnknownUserForbidden(t *testing.T) {
	s := newAppForTests(t, nil).server
	request := httptest.NewRequest("POST", "/sessions", strings.NewReader(`{"email":"something@soph.wiki","password":"pass"}`))
	assert.Equal(t, http.StatusForbidden, doRequest(s, request).Code)
}

func TestValidCredentialsAccepted(t *testing.T) {
	type response struct {
		Token string `json:"token"`
	}

	s := makeAccountWithPassword(t, "some-email@soph.wiki", "some-password")
	rr := doAuthenticate(s, `{"email":"some-email@soph.wiki","password":"some-password"}`)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
	assert.Equal(t, strconv.Itoa(rr.Body.Len()), rr.Header().Get("Content-Length"))

	var resp response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "auth="+resp.Token+"; SameSite=Strict; Secure; HttpOnly", rr.Header().Get("Set-Cookie"))
}

func TestInvalidCredentialsForbidden(t *testing.T) {
	s := makeAccountWithPassword(t, "some-email@soph.wiki", "some-password")
	rr := doAuthenticate(s, `{"email":"some-email@soph.wiki","password":"wrong-password"}`)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func makeAccountWithPassword(t *testing.T, email, password string) *endpoints.Server {
	app := newAppForTests(t, nil)
	rr := doSaveAccount(app.server, `{"email":"`+email+`"}`)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	emailInfo := app.emailer.welcomes[0]
	setPasswordBody := `{"password":"` + password + `","resetToken":"` + emailInfo.ResetToken + `"}`
	rr = doSetPassword(app.server, strconv.FormatInt(emailInfo.ID, 10), setPasswordBody)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	return app.server
}

func doAuthenticate(s *endpoints.Server, body string) *httptest.ResponseRecorder {
	return doRequest(s, httptest.NewRequest("POST", "/sessions", strings.NewReader(body)))
}

func assertMethodNotAllowed(t *testing.T, method string, path string) {
	t.Helper()
	s := newAppForTests(t, nil).server
	response := doRequest(s, httptest.NewRequest(method, path, nil))
	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
}
