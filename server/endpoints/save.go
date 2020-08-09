package endpoints

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/wikisophia/api/server/arguments"
)

// Implements POST /arguments
func saveHandler(saver arguments.Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request: "+err.Error(), http.StatusBadRequest)
			return
		}

		var arg arguments.Argument
		if err := json.Unmarshal(data, &arg); err != nil {
			http.Error(w, "Failed to unmarshal argument: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := arguments.ValidateArgument(arg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := saver.Save(context.Background(), arg)
		if err != nil {
			http.Error(w, "Failed to save argument: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		arg.ID = id
		arg.Version = 1
		w.Header().Set("Location", "/arguments/"+strconv.FormatInt(id, 10)+"/version/1")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		writeArgument(w, arg, strconv.FormatInt(id, 10))
	}
}
