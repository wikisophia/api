package endpoints_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveArgument(t *testing.T) {
	server, id, ok := newServerWithData(t, intendedOrigArg)
	if !ok {
		return
	}
	rr := doGetArgument(server, id)
	assertArgumentsMatch(t, intendedOrigArg, rr)
}

func TestSaveNoConclusion(t *testing.T) {
	rr := doSaveArgument(newServerForTests(), `{"premises":["Socrates is a man","All men are mortal"]}`)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveNoPremises(t *testing.T) {
	rr := doSaveArgument(newServerForTests(), `{"conclusion":"Socrates is mortal"}`)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveNotJSON(t *testing.T) {
	rr := doSaveArgument(newServerForTests(), `bad payload`)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
