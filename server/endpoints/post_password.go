package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api/server/accounts"
)

// Implements POST /accounts/:id/password
func setPasswordHandler(passwordSetter accounts.PasswordSetter) httprouter.Handle {
	type request struct {
		OldPassword string `json:"oldPassword"`
		Password    string `json:"password"`
		ResetToken  string `json:"resetToken"`
	}

	respondToStoreError := func(w http.ResponseWriter, err error) {
		if errors.As(err, &accounts.ProhibitedPasswordError{}) {
			http.Error(w, "Failed to set password: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errors.As(err, &accounts.InvalidResetTokenError{}) || errors.As(err, &accounts.InvalidPasswordError{}) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if errors.As(err, &accounts.AccountNotExistsError{}) {
			http.Error(w, "Unknown user ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
	}

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		idString := params.ByName("id")
		id, goodID := parseInt64Param(idString)
		if !goodID {
			http.Error(w, fmt.Sprintf("account %s does not exist", idString), http.StatusNotFound)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request: "+err.Error(), http.StatusBadRequest)
			return
		}

		var req request
		if err := json.Unmarshal(data, &req); err != nil {
			http.Error(w, "Malformed request: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.ResetToken != "" {
			if req.OldPassword != "" {
				http.Error(w, "Only one of oldPassword or resetToken should be defined.", http.StatusBadRequest)
				return
			}
			if err := passwordSetter.SetForgottenPassword(context.Background(), id, req.Password, req.ResetToken); err != nil {
				respondToStoreError(w, err)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if req.OldPassword == "" {
			http.Error(w, "Either oldPassword or resetToken must be defined.", http.StatusBadRequest)
			return
		}
		if err := passwordSetter.ChangePassword(context.Background(), id, req.OldPassword, req.Password); err != nil {
			respondToStoreError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
