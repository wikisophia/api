package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments"
)

func TestGetArgument(t *testing.T) {
	server := newServerForTests()
	id := doSaveObject(t, server, intendedOrigArg)
	rr := doGetArgument(server, id)
	if !assertSuccessfulJSON(t, rr) {
		return
	}
	actual := assertParseArgument(t, rr.Body.Bytes())
	assertArgumentsMatch(t, intendedOrigArg, actual)
}

func TestGetLatest(t *testing.T) {
	server := newServerForTests()
	id := doSaveObject(t, server, unintendedOrigArg)
	doUpdatePremises(t, server, id, updates)
	rr := doGetArgument(server, id)
	if !assertSuccessfulJSON(t, rr) {
		return
	}
	actual := assertParseArgument(t, rr.Body.Bytes())

	assertArgumentsMatch(t, arguments.Argument{
		Conclusion: unintendedOrigArg.Conclusion,
		Premises:   updates,
	}, actual)
}

func TestGetMissingArgument(t *testing.T) {
	rr := doGetArgument(newServerForTests(), 1)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetStringID(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
