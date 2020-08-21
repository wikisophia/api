package endpoints_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/arguments/argumentstest"
)

func TestPatchLive(t *testing.T) {
	original := argumentstest.ParseSample(t, "../samples/save-request.json")
	update := argumentstest.ParseSample(t, "../samples/update-request.json")

	server := newAppForTests(t, nil).server
	id := doSaveObject(t, server, original)
	update.ID = id

	rr := doValidUpdate(t, server, update)
	update.Version = 2
	responseBytes, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err)
	parsed := parseArgumentResponse(t, responseBytes)
	assert.Equal(t, update, parsed)

	rr = doGetArgument(server, id)
	assertSuccessfulJSON(t, rr)
	actual := parseArgumentResponse(t, rr.Body.Bytes())
	assert.Equal(t, update, actual)
}

func TestUpdateLocation(t *testing.T) {
	original := argumentstest.ParseSample(t, "../samples/save-request.json")
	update := argumentstest.ParseSample(t, "../samples/update-request.json")

	server := newAppForTests(t, nil).server
	id := doSaveObject(t, server, original)
	update.ID = id
	rr := doValidUpdate(t, server, update)
	assert.Equal(t, "/arguments/1/version/2", rr.Header().Get("Location"))
}

func TestPatchUnknown(t *testing.T) {
	server := newAppForTests(t, nil).server
	payload := string(argumentstest.ReadFile(t, "../samples/update-request.json"))
	rr := doRequest(server, httptest.NewRequest("PATCH", "/arguments/1", strings.NewReader(payload)))
	assert.Equal(t, http.StatusNotFound, rr.Code, "body: %s", rr.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestMalformedPatch(t *testing.T) {
	assertPatchRejected(t, "not json")
}

func TestPatchOnePremise(t *testing.T) {
	assertPatchRejected(t, `{"premises":["Socrates is a man"]}`)
}

func TestPatchEmpty(t *testing.T) {
	assertPatchRejected(t, `{"premises":["Socrates is a man", ""]}`)
}

func assertPatchRejected(t *testing.T, payload string) {
	t.Helper()
	original := argumentstest.ParseSample(t, "../samples/save-request.json")

	server := newAppForTests(t, nil).server
	id := doSaveObject(t, server, original)
	rr := doRequest(server, httptest.NewRequest("PATCH", "/arguments/"+strconv.FormatInt(id, 10), strings.NewReader(payload)))
	assert.Equal(t, http.StatusBadRequest, rr.Code, "body: %s", rr.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
