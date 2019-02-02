package endpoints_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/wikisophia/api-arguments/server/endpoints"

	"github.com/stretchr/testify/assert"

	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/config"
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
	assert.Equal(t, http.StatusBadRequest, rr.Code)
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

func assertSuccessfulJSON(t *testing.T, rr *httptest.ResponseRecorder) bool {
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	return !t.Failed()
}

func assertParseArgument(t *testing.T, data []byte) arguments.Argument {
	var argument arguments.Argument
	assert.NoError(t, json.Unmarshal(data, &argument))
	return argument
}

func assertParseAllArguments(t *testing.T, data []byte) endpoints.GetAllResponse {
	var getAll endpoints.GetAllResponse
	assert.NoError(t, json.Unmarshal(data, &getAll))
	return getAll
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

func doUpdatePremises(t *testing.T, server *endpoints.Server, id int64, updates ...[]string) {
	for _, update := range updates {
		arg := arguments.Argument{
			Premises: update,
		}
		updatePayload, err := json.Marshal(arg)
		if !assert.NoError(t, err) {
			return
		}
		rr := doPatchArgument(server, id, string(updatePayload))

		// Firefox parses empty response to AJAX calls as XML and throws an error.
		// The 204 response makes it works as expected.
		assert.Equal(t, http.StatusNoContent, rr.Code)
	}
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

func parseFile(t *testing.T, unixPath string, into interface{}) bool {
	fileBytes, err := ioutil.ReadFile(unixPath)
	if !assert.NoError(t, err) {
		return false
	}

	return assert.NoError(t, json.Unmarshal(fileBytes, into))
}

func newServerForTests() *endpoints.Server {
	cfg := config.Defaults()
	cfg.Storage.Type = config.StorageTypeMemory
	return endpoints.NewServer(cfg)
}
