package endpoints

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wikisophia/api/server/accounts"
)

type accountResetDependencies interface {
	accounts.ResetTokenGenerator
	accounts.Emailer
}

// Handle POST /accounts requests. This either registers a new account or
// generates a password reset token if the account already exists.
func accountHandler(dependencies accountResetDependencies) http.HandlerFunc {
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

		account, accountIsNew, err := dependencies.NewResetToken(context.Background(), req.Email)
		if err != nil {
			http.Error(w, "An internal error occurred. Please try again later.", http.StatusInternalServerError)
		}

		if accountIsNew {
			if err = dependencies.SendWelcome(context.Background(), account); err != nil {
				log.Printf("ERROR: Failed to send welcome email to %s: %v", account.Email, err)
			}
		} else {
			if err = dependencies.SendReset(context.Background(), account); err != nil {
				log.Printf("ERROR: Failed to send password reset email to %s: %v", account.Email, err)
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
