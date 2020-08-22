package http

import (
	"crypto/ecdsa"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api/server/accounts"
	"github.com/wikisophia/api/server/accounts/email"
)

type Dependencies interface {
	accounts.Store
	email.Emailer
}

// AppendRoutes populates the router with all the endpoints related to accounts.
func AppendRoutes(router *httprouter.Router, key *ecdsa.PrivateKey, dependencies Dependencies) {
	router.HandlerFunc("POST", "/accounts", accountHandler(dependencies))
	router.POST("/accounts/:id/password", setPasswordHandler(dependencies))
	router.HandlerFunc("POST", "/sessions", postSessionHandler(key, dependencies))
}
