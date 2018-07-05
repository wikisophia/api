package endpoints_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/wikisophia/api-arguments/endpoints"

	"github.com/stretchr/testify/assert"

	"github.com/wikisophia/api-arguments/arguments"
	"github.com/wikisophia/api-arguments/config"
)

var intendedOrigArg = arguments.Argument{
	Conclusion: "Socrates is mortal",
	Premises: []string{
		"Socrates is a man",
		"All men are mortal",
	},
}

var unintendedOrigArg = arguments.Argument{
	Conclusion: "Socrates is mortal",
	Premises: []string{
		"Socrates is a human",
		"All men are mortal",
	},
}

var updates = []string{
	"Socrates is a man",
	"All men are mortal",
}

func TestGetCollection(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments", nil))
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPostWithID(t *testing.T) {
	req := httptest.NewRequest("POST", "/arguments/1", nil)
	rr := doRequest(newServerForTests(), req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPostVersion(t *testing.T) {
	req := httptest.NewRequest("POST", "/arguments/1/version/1", nil)
	rr := doRequest(newServerForTests(), req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func assertArgumentsMatch(t *testing.T, expected arguments.Argument, rr *httptest.ResponseRecorder) {
	if !assert.Equal(t, http.StatusOK, rr.Code) {
		return
	}
	if !assert.Equal(t, "application/json", rr.Header().Get("Content-Type")) {
		return
	}
	var actual arguments.Argument
	if !assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &actual)) {
		return
	}

	assert.Equal(t, expected.Conclusion, actual.Conclusion)
	assert.ElementsMatch(t, expected.Premises, actual.Premises)
}

func newServerWithData(t *testing.T, argument arguments.Argument, updates ...[]string) (*endpoints.Server, int64, bool) {
	payload, err := json.Marshal(argument)
	if !assert.NoError(t, err) {
		return nil, -1, false
	}
	server := newServerForTests()
	rr := doSaveArgument(server, string(payload))
	if !assert.Equal(t, http.StatusCreated, rr.Code) {
		return nil, -1, false
	}
	id, err := parseArgumentID(rr.Header().Get("Location"))
	if !assert.NoError(t, err) {
		return nil, -1, false
	}
	for _, update := range updates {
		arg := arguments.Argument{
			Premises: update,
		}
		updatePayload, err := json.Marshal(arg)
		if !assert.NoError(t, err) {
			return nil, -1, false
		}
		rr := doPatchArgument(server, id, string(updatePayload))

		// Firefox parses empty response to AJAX calls as XML and throws an error.
		// The 204 response makes it works as expected.
		if !assert.Equal(t, http.StatusNoContent, rr.Code) {
			return nil, -1, false
		}
	}
	return server, id, true
}

func doGetArgument(server *endpoints.Server, id int64) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10), nil)
	return doRequest(server, req)
}

func doGetArgumentVersion(server *endpoints.Server, id int64, version int) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10)+"/version/"+strconv.Itoa(version), nil)
	return doRequest(server, req)
}

func doPatchArgument(server *endpoints.Server, id int64, payload string) *httptest.ResponseRecorder {
	return doRequest(server, httptest.NewRequest("PATCH", "/arguments/"+strconv.FormatInt(id, 10), strings.NewReader(payload)))
}

func doSaveArgument(server *endpoints.Server, payload string) *httptest.ResponseRecorder {
	return doRequest(server, httptest.NewRequest("POST", "/arguments", strings.NewReader(payload)))
}

func doRequest(server *endpoints.Server, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.Handle(rr, req)
	return rr
}

func parseArgumentID(location string) (int64, error) {
	if location == "" {
		return -1, errors.New("response Location header should not be empty")
	}
	capture := regexp.MustCompile(`/arguments/(.*)`)
	matches := capture.FindStringSubmatch(location)
	if len(matches) != 2 {
		return -1, fmt.Errorf("response Location header should be of the form /arguments/:id. Got %s", location)
	}
	idString := matches[1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		return -1, fmt.Errorf("response Location header /arguments/:id had an invalid ID: %s", idString)
	}
	return int64(id), nil
}

var configForTests = config.Configuration{
	Storage: config.Storage{
		Type: config.StorageTypeMemory,
	},
}

func newServerForTests() *endpoints.Server {
	return endpoints.NewServer(configForTests)
}
