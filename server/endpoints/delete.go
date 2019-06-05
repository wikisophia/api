package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ArgumentDeleter can delete arguments by ID.
type ArgumentDeleter interface {
	// Delete deletes an argument (and all its versions) from the site.
	// If the argument didn't exist, the error will be a NotFoundError.
	Delete(ctx context.Context, id int64) error
}

func deleteHandler(deleter ArgumentDeleter) httprouter.Handle {
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
