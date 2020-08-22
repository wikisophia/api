package http_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/acceptancetest"
	"github.com/wikisophia/api/server/arguments"
	argumentsHttp "github.com/wikisophia/api/server/arguments/http"
)

const samplesPath = "../../samples/"

func newApp(t *testing.T, cfg *acceptancetest.AppConfig) *app {
	return &app{
		App: acceptancetest.NewApp(t, cfg),
		t:   t,
	}
}

type app struct {
	*acceptancetest.App
	t *testing.T
}

func (a *app) GetLiveSuccessfully(id int64) arguments.Argument {
	rr := a.Do(httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10), nil))
	assert.Equal(a.t, http.StatusOK, rr.Code)
	assert.Equal(a.t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	var getOne argumentsHttp.GetOneResponse
	require.NoError(a.t, json.Unmarshal(rr.Body.Bytes(), &getOne))
	return getOne.Argument
}

func (a *app) GetVersionedSuccessfully(id int64, version int) arguments.Argument {
	rr := a.Do(httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10)+"/version/"+strconv.Itoa(version), nil))
	assert.Equal(a.t, http.StatusOK, rr.Code)
	assert.Equal(a.t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	var getOne argumentsHttp.GetOneResponse
	require.NoError(a.t, json.Unmarshal(rr.Body.Bytes(), &getOne))
	return getOne.Argument
}

func (a *app) FetchSome(options arguments.FetchSomeOptions) *httptest.ResponseRecorder {
	path := "/arguments"
	queryParamSeparator := newQueryParamSeparatorGenerator()

	if options.Conclusion != "" {
		path += queryParamSeparator() + "conclusion=" + url.QueryEscape(options.Conclusion)
	}
	if len(options.ConclusionContainsAll) > 0 {
		path += queryParamSeparator() + "search=" + url.QueryEscape(strings.Join(options.ConclusionContainsAll, " "))
	}
	if options.Count > 0 {
		path += queryParamSeparator() + "count=" + strconv.Itoa(options.Count)
	}
	if options.Offset > 0 {
		path += queryParamSeparator() + "offset=" + strconv.Itoa(options.Offset)
	}
	if len(options.Exclude) > 0 {
		s := make([]string, 0, len(options.Exclude))
		for i := 0; i < len(options.Exclude); i++ {
			s = append(s, strconv.FormatInt(options.Exclude[i], 10))
		}
		path += queryParamSeparator() + "exclude=" + strings.Join(s, "%2C")
	}
	req := httptest.NewRequest("GET", path, nil)
	return a.App.Do(req)
}

func (a *app) FetchSomeSuccessfully(t *testing.T, options arguments.FetchSomeOptions) []arguments.Argument {
	rr := a.FetchSome(options)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	var getAll argumentsHttp.GetAllResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &getAll))
	return getAll.Arguments
}

func (a *app) SaveSuccessfully(t *testing.T, argument arguments.Argument) int64 {
	payload, err := json.Marshal(argument)
	require.NoError(t, err)
	rr := a.App.Do(httptest.NewRequest("POST", "/arguments", bytes.NewReader(payload)))
	require.Equal(t, http.StatusCreated, rr.Code)
	id := parseArgumentID(t, rr.Header().Get("Location"))
	argument.ID = id
	argument.Version = 1
	responseBody := parseArgumentResponse(t, rr.Body.Bytes())
	assert.Equal(t, argument, responseBody)
	return id
}

func (a *app) SaveAllSuccessfully(t *testing.T, args []arguments.Argument) {
	for i := 0; i < len(args); i++ {
		id := a.SaveSuccessfully(t, args[i])
		args[i].ID = id
		args[i].Version = 1
	}
}

func (a *app) UpdateSuccessfully(t *testing.T, update arguments.Argument) (arguments.Argument, string) {
	id := update.ID
	update.ID = 0
	updatePayload, err := json.Marshal(update)
	assert.NoError(t, err)
	rr := a.Do(httptest.NewRequest("PATCH", "/arguments/"+strconv.FormatInt(id, 10), bytes.NewReader(updatePayload)))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
	responseBytes, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err)
	parsed := parseArgumentResponse(t, responseBytes)
	return parsed, rr.Header().Get("Location")
}

func newQueryParamSeparatorGenerator() func() string {
	hasQueryParam := false
	return func() string {
		if hasQueryParam {
			return "&"
		}
		hasQueryParam = true
		return "?"
	}
}

func parseArgumentResponse(t *testing.T, data []byte) arguments.Argument {
	var getOne argumentsHttp.GetOneResponse
	require.NoError(t, json.Unmarshal(data, &getOne))
	return getOne.Argument
}

func parseArgumentID(t *testing.T, location string) int64 {
	assert.NotEmpty(t, location)
	capture := regexp.MustCompile(`/arguments/(.*)/version/.*`)
	matches := capture.FindStringSubmatch(location)
	assert.Len(t, matches, 2)
	idString := matches[1]
	id, err := strconv.Atoi(idString)
	assert.NoError(t, err)
	return int64(id)
}

func assertArgumentSetsMatch(t *testing.T, expected []arguments.Argument, actual []arguments.Argument) {
	expectedMap := argumentListToMap(t, expected)
	actualMap := argumentListToMap(t, actual)
	assert.Equal(t, expectedMap, actualMap)
}

func argumentListToMap(t *testing.T, list []arguments.Argument) map[int64]arguments.Argument {
	theMap := make(map[int64]arguments.Argument)
	for i := 0; i < len(list); i++ {
		assert.NotContains(t, theMap, list[i].ID, "duplicate ID: %d", list[i].ID)
		theMap[list[i].ID] = list[i]
	}
	return theMap
}
