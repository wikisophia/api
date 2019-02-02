package endpoints_test

import (
	"testing"

	"github.com/wikisophia/api-arguments/server/endpoints"
)

func TestGetAll(t *testing.T) {
	var expected endpoints.GetAllResponse
	parseFile(t, "../samples/get-all-response.json", &expected)
	server := newServerForTests()
	id := doSaveObject(t, server, expected.Arguments[0].Argument)
	expected.Arguments[0].ID = id

	for i := 1; i < len(expected.Arguments); i++ {
		id := doSaveObject(t, server, expected.Arguments[i].Argument)
		expected.Arguments[i].ID = id
	}

	rr := doGetAllArguments(server, expected.Arguments[0].Conclusion)
	assertSuccessfulJSON(t, rr)
	actual := assertParseAllArguments(t, rr.Body.Bytes())
	assertArgumentSetsMatch(t, expected.Arguments, actual.Arguments)
}
