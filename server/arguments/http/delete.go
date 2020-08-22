package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api/server/arguments"
)

func deleteHandler(deleter arguments.Deleter) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseInt64Param(params.ByName("id"))
		if !goodID {
			http.Error(w, fmt.Sprintf("argument %s does not exist", params.ByName("id")), http.StatusNotFound)
			return
		}
		if err := deleter.Delete(context.Background(), id); writeStoreError(w, err) {
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}
}
