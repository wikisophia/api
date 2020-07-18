package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api-arguments/arguments"
)

// ArgumentUpdater can update existing arguments.
type ArgumentUpdater interface {
	// Update makes a new version of the argument. It returns the new argument's version.
	// If no argument with this ID exists, the returned error is an arguments.NotFoundError.
	Update(ctx context.Context, argument arguments.Argument) (version int, err error)
}

// Implements PATCH /arguments/:id
func updateHandler(updater ArgumentUpdater) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseInt64Param(params.ByName("id"))
		if !goodID || id < 1 {
			http.Error(w, fmt.Sprintf("argument %s does not exist", params.ByName("id")), http.StatusNotFound)
			return
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusInternalServerError)
			return
		}
		var arg arguments.Argument
		if err := json.Unmarshal(bodyBytes, &arg); err != nil {
			http.Error(w, "request body parse failure. Check the JSON syntax in your request body.", http.StatusBadRequest)
			return
		}
		if arg.ID != 0 {
			http.Error(w, "request.id should not be defined. The ID is taken from the URL path.", http.StatusBadRequest)
		}
		arg.ID = id
		if err := arguments.ValidateArgument(arg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		version, err := updater.Update(context.Background(), arg)
		if writeStoreError(w, err) {
			return
		}
		arg.Version = version

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Location", "/arguments/"+strconv.FormatInt(id, 10)+"/version/"+strconv.Itoa(int(version)))
		w.WriteHeader(http.StatusOK)
		writeArgument(w, arg, strconv.FormatInt(id, 10))
	}
}

func parseInt64Param(param string) (int64, bool) {
	parsed, err := strconv.ParseInt(param, 10, 0)
	return parsed, err == nil
}

func writeStoreError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(*arguments.NotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
		return true
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return true
}

func writeArgument(w http.ResponseWriter, arg arguments.Argument, id string) {
	data, err := json.Marshal(GetOneResponse{
		Argument: arg,
	})
	if err != nil {
		http.Error(w, "failed json.marshal on argument "+id, http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// GetOneResponse is the contract class for JSON responses of a single argument.
//
// Examples include:
//
//   GET /argument/{id}
//   GET /argument/{id}/version/{version}
//   POST /arguments
//   PATCH /argument?id=1
//
// etc.
//
type GetOneResponse struct {
	Argument arguments.Argument `json:"argument"`
}
