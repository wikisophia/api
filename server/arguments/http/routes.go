package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api/server/arguments"
)

// AppendRoutes populates the router with all the /arguments* endpoints.
func AppendRoutes(router *httprouter.Router, store arguments.Store) {
	router.HandlerFunc("POST", "/arguments", saveHandler(store))
	router.HandlerFunc("GET", "/arguments", getAllArgumentsHandler(store))
	router.GET("/arguments/:id", getLiveArgumentHandler(store))
	router.PATCH("/arguments/:id", updateHandler(store))
	router.DELETE("/arguments/:id", deleteHandler(store))
	router.GET("/arguments/:id/version/:version", getArgumentByVersionHandler(store))
}
