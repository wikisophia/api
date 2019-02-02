package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	server := newServerForTests()
	id := doSaveObject(t, server, unintendedOrigArg)
	doUpdatePremises(t, server, id, updates)
	rr := doGetArgumentVersion(server, id, 1)
	if !assertSuccessfulJSON(t, rr) {
		return
	}
	actual := assertParseArgument(t, rr.Body.Bytes())
	assertArgumentsMatch(t, unintendedOrigArg, actual)
}

func TestGetMissingVersion(t *testing.T) {
	server := newServerForTests()
	id := doSaveObject(t, server, unintendedOrigArg)
	rr := doGetArgumentVersion(server, id, 100)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetStringVersion(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments/1/version/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetLargeVersion(t *testing.T) {
	rr := doGetArgumentVersion(newServerForTests(), 1, 65537)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
