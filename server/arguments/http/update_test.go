package http_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestPatchLive(t *testing.T) {
	original := acceptancetest.ParseSample(t, samplesPath+"save-request.json")
	update := acceptancetest.ParseSample(t, samplesPath+"update-request.json")

	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, original)
	update.ID = id

	parsed, _ := app.UpdateSuccessfully(t, update)
	update.Version = 2
	assert.Equal(t, update, parsed)

	actual := app.GetLiveSuccessfully(id)
	assert.Equal(t, update, actual)
}

func TestUpdateLocation(t *testing.T) {
	original := acceptancetest.ParseSample(t, samplesPath+"save-request.json")
	update := acceptancetest.ParseSample(t, samplesPath+"update-request.json")

	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, original)
	update.ID = id
	_, location := app.UpdateSuccessfully(t, update)
	assert.Equal(t, "/arguments/1/version/2", location)
}

func TestPatchUnknown(t *testing.T) {
	app := newApp(t, nil)
	payload := string(acceptancetest.ReadFile(t, samplesPath+"update-request.json"))
	rr := app.Do(httptest.NewRequest("PATCH", "/arguments/1", strings.NewReader(payload)))
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
	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, acceptancetest.ParseSample(t, samplesPath+"save-request.json"))
	rr := app.Do(httptest.NewRequest("PATCH", "/arguments/"+strconv.FormatInt(id, 10), strings.NewReader(payload)))
	assert.Equal(t, http.StatusBadRequest, rr.Code, "body: %s", rr.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
