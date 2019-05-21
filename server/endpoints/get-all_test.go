package endpoints_test

import (
	"testing"

	"github.com/wikisophia/api-arguments/server/arguments/argumentstest"
)

func TestGetAll(t *testing.T) {
	expected := parseGetAllResponse(t, argumentstest.ReadFile(t, "../samples/get-all-response.json"))
	server := newServerForTests()
	id := doSaveObject(t, server, expected.Arguments[0])
	expected.Arguments[0].ID = id
	expected.Arguments[0].Version = 1

	for i := 1; i < len(expected.Arguments); i++ {
		id := doSaveObject(t, server, expected.Arguments[i])
		expected.Arguments[i].ID = id
		expected.Arguments[i].Version = 1
	}

	rr := doGetAllArguments(server, expected.Arguments[0].Conclusion)
	assertSuccessfulJSON(t, rr)
	actual := parseGetAllResponse(t, rr.Body.Bytes())
	assertArgumentSetsMatch(t, expected.Arguments, actual.Arguments)
}
