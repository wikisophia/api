package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api-arguments/server/arguments"
)

// Implements PATCH /arguments/:id
func (s *Server) updateArgument() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseIntParam(params.ByName("id"))
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

		version, err := s.argumentStore.Update(context.Background(), arg)
		if writeStoreError(w, err) {
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Location", "/arguments/"+strconv.FormatInt(id, 10)+"/version/"+strconv.Itoa(int(version)))
		w.WriteHeader(http.StatusNoContent)
	}
}

func parseIntParam(param string) (int64, bool) {
	parsed, err := strconv.ParseInt(param, 10, 0)
	return parsed, err == nil
}

// shrink an int to 16 bits. Return true if it still holds the same value
func shrinkInt16(value int64) (int16, bool) {
	shrunk := int16(value)
	return shrunk, int64(shrunk) == value
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
	data, err := json.Marshal(arg)
	if err != nil {
		http.Error(w, "failed json.marshal on argument "+id, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
