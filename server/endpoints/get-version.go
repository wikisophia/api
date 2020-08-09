package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api/server/arguments"
)

// Implements GET /arguments/:id/version/:version
func getArgumentByVersionHandler(getter arguments.GetVersioned) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseInt64Param(params.ByName("id"))
		version, ok := parseIntParam(params.ByName("version"))

		if !goodID || !ok {
			response := fmt.Sprintf("version %s of argument %s does not exist", params.ByName("version"), params.ByName("id"))
			http.Error(w, response, http.StatusNotFound)
			return
		}
		arg, err := getter.FetchVersion(context.Background(), id, version)
		if writeStoreError(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		writeArgument(w, arg, params.ByName("id"))
	}
}
