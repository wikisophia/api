package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments"
)

func TestPatchLive(t *testing.T) {
	server := newServerForTests()
	id := doSaveObject(t, server, intendedOrigArg)
	doUpdatePremises(t, server, id, updates)
	rr := doGetArgument(server, id)
	if !assertSuccessfulJSON(t, rr) {
		return
	}
	actual := assertParseArgument(t, rr.Body.Bytes())
	assertArgumentsMatch(t, arguments.Argument{
		Conclusion: unintendedOrigArg.Conclusion,
		Premises:   intendedOrigArg.Premises,
	}, actual)
}

func TestPatchUnknown(t *testing.T) {
	server := newServerForTests()
	payload := `{"premises":["Socrates is a man", "All men are mortal"]}`
	rr := doRequest(server, httptest.NewRequest("PATCH", "/arguments/1", strings.NewReader(payload)))
	assert.Equal(t, http.StatusNotFound, rr.Code, "body: %s", rr.Body.String())
}

func TestMalformedPatch(t *testing.T) {
	assertPatchRejected(t, "not json")
}

func TestPatchConclusion(t *testing.T) {
	assertPatchRejected(t, `{"conclusion":"Socrates is mortal","premises":["Socrates is a man", "All men are mortal"]}`)
}

func TestPatchOnePremise(t *testing.T) {
	assertPatchRejected(t, `{"premises":["Socrates is a man"]}`)
}

func TestPatchEmpty(t *testing.T) {
	assertPatchRejected(t, `{"premises":["Socrates is a man", ""]}`)
}

func assertPatchRejected(t *testing.T, payload string) {
	t.Helper()
	server := newServerForTests()
	id := doSaveObject(t, server, unintendedOrigArg)
	rr := doRequest(server, httptest.NewRequest("PATCH", "/arguments/"+strconv.FormatInt(id, 10), strings.NewReader(payload)))
	assert.Equal(t, http.StatusBadRequest, rr.Code, "body: %s", rr.Body.String())
}
