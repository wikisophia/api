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
	server, id, ok := newServerWithData(t, unintendedOrigArg, updates)
	if !ok {
		return
	}
	rr := doGetArgument(server, id)
	assertArgumentsMatch(t, arguments.Argument{
		Conclusion: unintendedOrigArg.Conclusion,
		Premises:   intendedOrigArg.Premises,
	}, rr)
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
	server, id, ok := newServerWithData(t, unintendedOrigArg)
	if !ok {
		return
	}
	rr := doRequest(server, httptest.NewRequest("PATCH", "/arguments/"+strconv.FormatInt(id, 10), strings.NewReader(payload)))
	assert.Equal(t, http.StatusBadRequest, rr.Code, "body: %s", rr.Body.String())
}
