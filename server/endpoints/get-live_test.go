package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments"
)

func TestGetArgument(t *testing.T) {
	server, id, ok := newServerWithData(t, intendedOrigArg)
	if !ok {
		return
	}
	rr := doGetArgument(server, id)
	assertArgumentsMatch(t, intendedOrigArg, rr)
}

func TestGetLatest(t *testing.T) {
	server, id, ok := newServerWithData(t, unintendedOrigArg, updates)
	if !ok {
		return
	}
	rr := doGetArgument(server, id)
	assertArgumentsMatch(t, arguments.Argument{
		Conclusion: unintendedOrigArg.Conclusion,
		Premises:   updates,
	}, rr)
}

func TestGetMissingArgument(t *testing.T) {
	rr := doGetArgument(newServerForTests(), 1)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetStringID(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
