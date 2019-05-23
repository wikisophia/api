package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api-arguments/server/arguments"
)

// ArgumentGetterLiveVersion can fetch the live version of an argument.
type ArgumentGetterLiveVersion interface {
	// FetchLive should return the latest "active" version of an argument.
	// If no argument with this ID exists, the error should be an arguments.NotFoundError.
	FetchLive(ctx context.Context, id int64) (arguments.Argument, error)
}

// Implements GET /arguments/:id
func getLiveArgumentHandler(getter ArgumentGetterLiveVersion) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseIntParam(params.ByName("id"))
		if !goodID {
			http.Error(w, fmt.Sprintf("argument %s does not exist", params.ByName("id")), http.StatusNotFound)
			return
		}

		arg, err := getter.FetchLive(context.Background(), id)
		if writeStoreError(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		writeArgument(w, arg, params.ByName("id"))
	}
}
