package endpoints_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/endpoints"
)

// This file has a bunch of helper methods used throughout the test code
// in this package.

// newServerForTests returns a Server that stores arguments in memory.
func newServerForTests() *endpoints.Server {
	cfg := config.Defaults()
	cfg.Storage.Type = config.StorageTypeMemory
	return endpoints.NewServer(cfg)
}

func readFile(t *testing.T, unixPath string) []byte {
	fileBytes, err := ioutil.ReadFile(filepath.FromSlash(unixPath))
	assert.NoError(t, err)
	return fileBytes
}

func parseArgument(t *testing.T, data []byte) arguments.Argument {
	var argument arguments.Argument
	assert.NoError(t, json.Unmarshal(data, &argument))
	return argument
}

func parseGetAllResponse(t *testing.T, data []byte) endpoints.GetAllResponse {
	var getAll endpoints.GetAllResponse
	assert.NoError(t, json.Unmarshal(data, &getAll))
	return getAll
}

func parseArgumentID(t *testing.T, location string) int64 {
	assert.NotEmpty(t, location)
	capture := regexp.MustCompile(`/arguments/(.*)`)
	matches := capture.FindStringSubmatch(location)
	assert.Len(t, matches, 2)
	idString := matches[1]
	id, err := strconv.Atoi(idString)
	assert.NoError(t, err)
	return int64(id)
}

func parseFile(t *testing.T, unixPath string, into interface{}) {
	fileBytes, err := ioutil.ReadFile(filepath.FromSlash(unixPath))
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(fileBytes, into))
}

func assertSuccessfulJSON(t *testing.T, rr *httptest.ResponseRecorder) bool {
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
	return !t.Failed()
}

func assertArgumentsMatch(t *testing.T, expected arguments.Argument, actual arguments.Argument) {
	assert.Equal(t, expected.Conclusion, actual.Conclusion)
	assert.ElementsMatch(t, expected.Premises, actual.Premises)
}

func assertArgumentSetsMatch(t *testing.T, expected []arguments.ArgumentWithID, actual []arguments.ArgumentWithID) {
	expectedMap := argumentListToMap(t, expected)
	actualMap := argumentListToMap(t, actual)
	assert.Equal(t, expectedMap, actualMap)
}

func argumentListToMap(t *testing.T, list []arguments.ArgumentWithID) map[int64]arguments.Argument {
	theMap := make(map[int64]arguments.Argument)
	for i := 0; i < len(list); i++ {
		assert.NotContains(t, theMap, list[i].ID, "duplicate ID: %d", list[i].ID)
		theMap[list[i].ID] = list[i].Argument
	}
	return theMap
}

func doSaveObject(t *testing.T, server *endpoints.Server, argument arguments.Argument) int64 {
	payload, err := json.Marshal(argument)
	if !assert.NoError(t, err) {
		return -1
	}
	rr := doSaveArgument(server, string(payload))
	if !assert.Equal(t, http.StatusCreated, rr.Code) {
		return -1
	}
	id := parseArgumentID(t, rr.Header().Get("Location"))
	if !assert.NoError(t, err) {
		return -1
	}
	return id
}

func doValidUpdate(t *testing.T, server *endpoints.Server, id int64, update []string) *httptest.ResponseRecorder {
	wrapper := arguments.Argument{
		Premises: update,
	}
	updatePayload, err := json.Marshal(wrapper)
	assert.NoError(t, err)
	rr := doPatchArgument(server, id, string(updatePayload))

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
	return rr
}

func doGetArgument(server *endpoints.Server, id int64) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10), nil)
	return doRequest(server, req)
}

func doGetAllArguments(server *endpoints.Server, conclusion string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/arguments?conclusion="+url.QueryEscape(conclusion), nil)
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
