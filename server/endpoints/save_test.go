package endpoints_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments/argumentstest"
)

func TestSaveGetRoundtrip(t *testing.T) {
	expected := argumentstest.ParseSample(t, "../samples/save-request.json")
	server := newServerForTests()
	id := doSaveObject(t, server, expected)
	expected.ID = id
	expected.Version = 1
	rr := doGetArgument(server, id)
	assertSuccessfulJSON(t, rr)
	actual := argumentstest.ParseJSON(t, rr.Body.Bytes())
	assert.Equal(t, expected, actual)
}

func TestSaveNoConclusion(t *testing.T) {
	rr := doSaveArgument(newServerForTests(), `{"premises":["Socrates is a man","All men are mortal"]}`)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestSaveNoPremises(t *testing.T) {
	rr := doSaveArgument(newServerForTests(), `{"conclusion":"Socrates is mortal"}`)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestSaveNotJSON(t *testing.T) {
	rr := doSaveArgument(newServerForTests(), `bad payload`)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
