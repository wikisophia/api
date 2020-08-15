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

// TODO: When emails go out find a way to mock it and check values during tests.

func TestSetPasswordRejectsBadRequestsProperly(t *testing.T) {
	s := newServerForTests()
	rr := doSaveAccount(s, `{"email":"some-email@soph.wiki"}`)
	require.Equal(t, http.StatusNoContent, rr.Code)

	assert.Equal(t, http.StatusNotFound, doSetPassword(s, "2", `{"password":"password","resetToken":"abc"}`).Code)
	assert.Equal(t, http.StatusNotFound, doSetPassword(s, "non-numeric", "").Code)
	assert.Equal(t, http.StatusUnauthorized, doSetPassword(s, "1", `{"password":"abc","resetToken":"*"}`).Code)
	assertBadRequest(t, doSetPassword(s, "1", "not json"))
	assertBadRequest(t, doSetPassword(s, "1", "5"))
	assertBadRequest(t, doSetPassword(s, "1", "true"))
	assertBadRequest(t, doSetPassword(s, "1", "null"))
	assertBadRequest(t, doSetPassword(s, "1", "\"\""))
	assertBadRequest(t, doSetPassword(s, "1", "{}"))
	assertBadRequest(t, doSetPassword(s, "1", `{"password":"something"}`))
	assertBadRequest(t, doSetPassword(s, "1", `{"password":"something","resetToken":"abc","oldPassword":"something-else"}`))
}

func doSetPassword(s *endpoints.Server, id string, body string) *httptest.ResponseRecorder {
	return doRequest(s, httptest.NewRequest("POST", "/accounts/"+id+"/password", strings.NewReader(body)))
}

func assertNotFound(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
