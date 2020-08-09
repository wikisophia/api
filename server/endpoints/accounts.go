package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ResetTokenGenerator creates accounts and generate password reset tokens
type ResetTokenGenerator interface {
	// NewResetTokenWithAccount assigns a new password reset token to the account
	// with this email. If no accounts exist with this email, one will be created.
	NewResetTokenWithAccount(email string) (string, error)
}

// Handle POST /accounts requests. This either registers a new account or
// generates a password reset token if the account already exists.
func accountHandler(tokenGenerator ResetTokenGenerator) http.HandlerFunc {
	type request struct {
		Email string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var req request
		if err = json.Unmarshal(payload, &req); err != nil {
			http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Email == "" {
			http.Error(w, "missing required property: \"email\"", http.StatusBadRequest)
			return
		}

		_, err = tokenGenerator.NewResetTokenWithAccount(req.Email)
		if err != nil {
			http.Error(w, "An internal error occurred. Please try again later.", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
