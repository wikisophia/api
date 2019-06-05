package endpoints_test

import (
	"testing"

	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/arguments/argumentstest"
	"github.com/wikisophia/api-arguments/server/endpoints"
)

func TestGetAll(t *testing.T) {
	expected := parseGetAllResponse(t, argumentstest.ReadFile(t, "../samples/get-all-response.json"))
	server := newServerForTests()
	addAllArguments(t, server, expected.Arguments)
	assertFetchSome(t, server, arguments.FetchSomeOptions{}, expected.Arguments)
}

func TestGetWithConclusion(t *testing.T) {
	expected := parseGetAllResponse(t, argumentstest.ReadFile(t, "../samples/get-all-response.json"))
	server := newServerForTests()
	addAllArguments(t, server, expected.Arguments)

	for i := 1; i < len(expected.Arguments); i++ {
		if expected.Arguments[i].Conclusion != expected.Arguments[0].Conclusion {
			expected.Arguments = append(expected.Arguments[:i], expected.Arguments[i+1:]...)
			i--
		}
	}
	assertFetchSome(t, server, arguments.FetchSomeOptions{
		Conclusion: expected.Arguments[0].Conclusion,
	}, expected.Arguments)
}

func TestGetWithOffsets(t *testing.T) {
	expected := parseGetAllResponse(t, argumentstest.ReadFile(t, "../samples/get-all-response.json"))
	server := newServerForTests()
	addAllArguments(t, server, expected.Arguments)
	assertFetchSome(t, server, arguments.FetchSomeOptions{
		Count: 1,
	}, []arguments.Argument{expected.Arguments[0]})
	assertFetchSome(t, server, arguments.FetchSomeOptions{
		Count:  1,
		Offset: 1,
	}, []arguments.Argument{expected.Arguments[1]})
}

func TestGetWithExclusions(t *testing.T) {
	expected := parseGetAllResponse(t, argumentstest.ReadFile(t, "../samples/get-all-response.json"))
	server := newServerForTests()
	addAllArguments(t, server, expected.Arguments)
	assertFetchSome(t, server, arguments.FetchSomeOptions{
		Exclude: []int64{expected.Arguments[0].ID},
	}, expected.Arguments[1:])
	assertFetchSome(t, server, arguments.FetchSomeOptions{
		Exclude: []int64{expected.Arguments[1].ID},
	}, append(append([]arguments.Argument{}, expected.Arguments[0]), expected.Arguments[2:]...))
}

func addAllArguments(t *testing.T, server *endpoints.Server, args []arguments.Argument) {
	for i := 0; i < len(args); i++ {
		id := doSaveObject(t, server, args[i])
		args[i].ID = id
		args[i].Version = 1
	}
}

func assertFetchSome(t *testing.T, server *endpoints.Server, options arguments.FetchSomeOptions, expected []arguments.Argument) {
	rr := doFetchSomeArguments(server, options)
	assertSuccessfulJSON(t, rr)
	actual := parseGetAllResponse(t, rr.Body.Bytes())
	assertArgumentSetsMatch(t, expected, actual.Arguments)
}
