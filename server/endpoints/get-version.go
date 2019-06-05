package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api-arguments/server/arguments"
)

// ArgumentGetterByVersion returns a specific version of an argument.
type ArgumentGetterByVersion interface {
	// FetchVersion should return a particular version of an argument.
	// If the the argument didn't exist, the error should be an arguments.NotFoundError.
	FetchVersion(ctx context.Context, id int64, version int) (arguments.Argument, error)
}

// Implements GET /arguments/:id/version/:version
func getArgumentByVersionHandler(getter ArgumentGetterByVersion) httprouter.Handle {
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
