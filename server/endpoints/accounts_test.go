package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/endpoints"
)

func TestAccountAcceptsEmails(t *testing.T) {
	rr := doSaveAccount(newAppForTests(testServerConfig{}).server, `{"email":"some-email@soph.wiki"}`)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestAccountRejectsBadRequestBodies(t *testing.T) {
	assertBadRequest(t, doSaveAccount(newAppForTests(testServerConfig{}).server, "not json"))
	assertBadRequest(t, doSaveAccount(newAppForTests(testServerConfig{}).server, "{}"))
	assertBadRequest(t, doSaveAccount(newAppForTests(testServerConfig{}).server, `{"email":null}`))
	assertBadRequest(t, doSaveAccount(newAppForTests(testServerConfig{}).server, `{"email":5}`))
	assertBadRequest(t, doSaveAccount(newAppForTests(testServerConfig{}).server, `{"email":true}`))
	assertBadRequest(t, doSaveAccount(newAppForTests(testServerConfig{}).server, `{"email":3.4}`))
}

func doSaveAccount(s *endpoints.Server, body string) *httptest.ResponseRecorder {
	return doRequest(s, httptest.NewRequest("POST", "/accounts", strings.NewReader(body)))
}

func assertBadRequest(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
