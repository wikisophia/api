package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/acceptancetest"
	"github.com/wikisophia/api/server/accounts"
)

func newApp(t *testing.T, cfg *acceptancetest.AppConfig) *app {
	return &app{
		App: acceptancetest.NewApp(t, cfg),
		t:   t,
	}
}

type app struct {
	*acceptancetest.App
	t *testing.T
}

func (a *app) SaveAccount(email string) *httptest.ResponseRecorder {
	type request struct {
		Email string `json:"email"`
	}
	data, err := json.Marshal(request{email})
	require.NoError(a.t, err)
	return a.Do(httptest.NewRequest("POST", "/accounts", bytes.NewReader(data)))
}

func (a *app) SaveAccountSuccessfully(email string) accounts.Account {
	numAccounts := len(a.Emailer.Welcomes)
	require.Equal(a.t, http.StatusNoContent, a.SaveAccount(email).Code)
	require.Len(a.t, a.Emailer.Welcomes, numAccounts+1)
	return *a.Emailer.Welcomes[numAccounts]
}

func (a *app) SaveAccountWithPasswordSuccessfully(email, password string) {
	account := a.SaveAccountSuccessfully(email)
	assert.Equal(a.t, http.StatusNoContent, a.ResetPassword(account.ID, account.ResetToken, password).Code)
}

func (a *app) ResetPassword(id int64, resetToken, password string) *httptest.ResponseRecorder {
	type request struct {
		Password   string `json:"password"`
		ResetToken string `json:"resetToken"`
	}
	data, err := json.Marshal(request{password, resetToken})
	require.NoError(a.t, err)
	return a.Do(httptest.NewRequest("POST", "/accounts/"+strconv.FormatInt(id, 10)+"/password", bytes.NewReader(data)))
}

func (a *app) ResetPasswordSuccessfully(id int64, resetToken, password string) {
	require.Equal(a.t, http.StatusNoContent, a.ResetPassword(id, resetToken, password).Code)
}

func (a *app) Authenticate(email, password string) *httptest.ResponseRecorder {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	data, err := json.Marshal(request{email, password})
	require.NoError(a.t, err)
	return a.Do(httptest.NewRequest("POST", "/sessions", bytes.NewReader(data)))
}

func (a *app) AuthenticateSuccessfully(email, password string) string {
	type response struct {
		Token string `json:"token"`
	}

	rr := a.Authenticate(email, password)
	require.Equal(a.t, http.StatusOK, rr.Code)
	require.Equal(a.t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
	require.Equal(a.t, strconv.Itoa(rr.Body.Len()), rr.Header().Get("Content-Length"))
	var r response
	require.NoError(a.t, json.Unmarshal(rr.Body.Bytes(), &r))
	assert.Equal(a.t, "auth="+r.Token+"; SameSite=Strict; Secure; HttpOnly", rr.Header().Get("Set-Cookie"))
	return r.Token
}

func (a *app) UpdatePassword(id int64, oldPassword, newPassword string) *httptest.ResponseRecorder {
	type request struct {
		OldPassword string `json:"oldPassword"`
		Password    string `json:"password"`
	}
	data, err := json.Marshal(request{oldPassword, newPassword})
	require.NoError(a.t, err)
	return a.Do(httptest.NewRequest("POST", "/accounts/"+strconv.FormatInt(id, 10)+"/password", bytes.NewReader(data)))
}

func (a *app) UpdatePasswordSuccessfully(id int64, oldPassword, newPassword string) {
	rr := a.UpdatePassword(id, oldPassword, newPassword)
	require.Equal(a.t, http.StatusNoContent, rr.Code)
}
