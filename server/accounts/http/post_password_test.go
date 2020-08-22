package http_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordSetsProperly(t *testing.T) {
	a := newApp(t, nil)
	acct := a.SaveAccountSuccessfully("some-email@soph.wiki")
	require.Equal(t, http.StatusForbidden, a.ResetPassword(acct.ID, "wrong-"+acct.ResetToken, "some-password").Code)
	a.ResetPasswordSuccessfully(acct.ID, acct.ResetToken, "some-password")
	a.AuthenticateSuccessfully("some-email@soph.wiki", "some-password")

	require.Equal(t, http.StatusForbidden, a.UpdatePassword(acct.ID, "some-wrong-password", "some-new-password").Code)
	a.UpdatePasswordSuccessfully(acct.ID, "some-password", "some-new-password")

	require.Equal(t, http.StatusForbidden, a.Authenticate("some-email@soph.wiki", "some-password").Code)
	a.AuthenticateSuccessfully("some-email@soph.wiki", "some-new-password")
}

func TestSetPasswordRejectsBadRequestsProperly(t *testing.T) {
	a := newApp(t, nil)
	validReset := `{"resetToken":"abc","password":"something"}`
	assert.Equal(t, http.StatusForbidden,
		a.Do(httptest.NewRequest("POST", "/accounts/1/password", strings.NewReader(validReset))).Code)
	assert.Equal(t, http.StatusForbidden,
		a.Do(httptest.NewRequest("POST", "/accounts/non-numeric/password", strings.NewReader(validReset))).Code)

	idString := strconv.FormatInt(a.SaveAccountSuccessfully("some-email@soph.wiki").ID, 10)
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", "not json")
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", "5")
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", "true")
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", "null")
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", `""`)
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", "{}")
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", `{"password":"something"}`)
	a.AssertBadRequest("POST", "/accounts/"+idString+"/password", `{"password":"something","resetToken":"abc","oldPassword":"something-else"}`)
}
