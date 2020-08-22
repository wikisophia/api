package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestSaveGetRoundtrip(t *testing.T) {
	expected := acceptancetest.ParseSample(t, samplesPath+"save-request.json")
	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, expected)
	expected.ID = id
	expected.Version = 1
	actual := app.GetLiveSuccessfully(id)
	assert.Equal(t, expected, actual)
}

func TestSaveNoConclusion(t *testing.T) {
	rr := newApp(t, nil).Do(newPostArgument(`{"premises":["Socrates is a man","All men are mortal"]}`))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestSaveNoPremises(t *testing.T) {
	rr := newApp(t, nil).Do(newPostArgument(`{"conclusion":"Socrates is mortal"}`))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestSaveNotJSON(t *testing.T) {
	rr := newApp(t, nil).Do(newPostArgument("bad payload"))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func newPostArgument(payload string) *http.Request {
	return httptest.NewRequest("POST", "/arguments", strings.NewReader(payload))
}
