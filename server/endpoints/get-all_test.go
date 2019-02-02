package endpoints_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/endpoints"
)

func TestGetAll(t *testing.T) {
	var expected endpoints.GetAllResponse
	if !parseFile(t, "../samples/get-all-response.json", &expected) {
		return
	}
	server := newServerForTests()
	id := doSaveObject(t, server, expected.Arguments[0].Argument)
	expected.Arguments[0].ID = id

	for i := 1; i < len(expected.Arguments); i++ {
		saved := doSaveObject(t, server, expected.Arguments[i].Argument)
		expected.Arguments[i].ID = saved
	}

	rr := doGetAllArguments(server, expected.Arguments[0].Conclusion)
	assertSuccessfulJSON(t, rr)
	actual := assertParseAllArguments(t, rr.Body.Bytes())
	assert.Equal(t, expected, actual)
}

func parseFile(t *testing.T, unixPath string, into interface{}) bool {
	fileBytes, err := ioutil.ReadFile(unixPath)
	if !assert.NoError(t, err) {
		return false
	}

	return assert.NoError(t, json.Unmarshal(fileBytes, into))
}
