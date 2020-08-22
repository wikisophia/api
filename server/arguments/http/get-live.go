package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api/server/arguments"
)

// Implements GET /arguments/:id
func getLiveArgumentHandler(getter arguments.GetLive) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseInt64Param(params.ByName("id"))
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
