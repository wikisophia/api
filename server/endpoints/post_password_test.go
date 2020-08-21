package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/endpoints"
)

func TestPasswordSetsProperly(t *testing.T) {
	app := newAppForTests(t, nil)
	require.Equal(t, http.StatusNoContent, doSaveAccount(app.server, `{"email":"some-email@soph.wiki"}`).Code)
	welcomeEmail := app.emailer.welcomes[0]
	rr := doSetPassword(app.server, strconv.FormatInt(welcomeEmail.ID, 10), `{"password":"some-password","resetToken":"wrong-`+welcomeEmail.ResetToken+`"}`)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
	rr = doSetPassword(app.server, strconv.FormatInt(welcomeEmail.ID, 10), `{"password":"some-password","resetToken":"`+welcomeEmail.ResetToken+`"}`)
	require.Equal(t, http.StatusNoContent, rr.Code)

	// TODO: Authenticate once the endpoint exists and make sure the password sets properly

	rr = doSetPassword(app.server, strconv.FormatInt(welcomeEmail.ID, 10), `{"password":"some-new-password","oldPassword":"some-wrong-password"}`)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
	rr = doSetPassword(app.server, strconv.FormatInt(welcomeEmail.ID, 10), `{"password":"some-new-password","oldPassword":"some-password"}`)
	require.Equal(t, http.StatusNoContent, rr.Code)

	// TODO: Authenticate again once that endpoint exists to make sure the password changes properly
}

func TestSetPasswordRejectsBadRequestsProperly(t *testing.T) {
	app := newAppForTests(t, nil)
	s := app.server
	rr := doSaveAccount(s, `{"email":"some-email@soph.wiki"}`)
	require.Equal(t, http.StatusNoContent, rr.Code)
	id := app.emailer.welcomes[0].ID
	idString := strconv.FormatInt(id, 10)

	assert.Equal(t, http.StatusNotFound, doSetPassword(s, strconv.FormatInt(id+1, 10), `{"password":"password","resetToken":"abc"}`).Code)
	assert.Equal(t, http.StatusNotFound, doSetPassword(s, "non-numeric", "").Code)
	assert.Equal(t, http.StatusUnauthorized, doSetPassword(s, idString, `{"password":"abc","resetToken":"*"}`).Code)
	assertBadRequest(t, doSetPassword(s, idString, "not json"))
	assertBadRequest(t, doSetPassword(s, idString, "5"))
	assertBadRequest(t, doSetPassword(s, idString, "true"))
	assertBadRequest(t, doSetPassword(s, idString, "null"))
	assertBadRequest(t, doSetPassword(s, idString, "\"\""))
	assertBadRequest(t, doSetPassword(s, idString, "{}"))
	assertBadRequest(t, doSetPassword(s, idString, `{"password":"something"}`))
	assertBadRequest(t, doSetPassword(s, idString, `{"password":"something","resetToken":"abc","oldPassword":"something-else"}`))
}

func doSetPassword(s *endpoints.Server, id string, body string) *httptest.ResponseRecorder {
	return doRequest(s, httptest.NewRequest("POST", "/accounts/"+id+"/password", strings.NewReader(body)))
}

func assertNotFound(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
