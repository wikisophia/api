package http_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/acceptancetest"
	"github.com/wikisophia/api/server/arguments"
	argumentsHttp "github.com/wikisophia/api/server/arguments/http"
)

func TestGetAll(t *testing.T) {
	expected := parseGetAllResponse(t, acceptancetest.ReadFile(t, samplesPath+"get-all-response.json"))
	a := newApp(t, nil)
	a.SaveAllSuccessfully(t, expected.Arguments)
	fetched := a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{})
	assertArgumentSetsMatch(t, expected.Arguments, fetched)
}

func TestGetWithConclusion(t *testing.T) {
	expected := parseGetAllResponse(t, acceptancetest.ReadFile(t, samplesPath+"get-all-response.json"))
	a := newApp(t, nil)
	a.SaveAllSuccessfully(t, expected.Arguments)

	for i := 1; i < len(expected.Arguments); i++ {
		if expected.Arguments[i].Conclusion != expected.Arguments[0].Conclusion {
			expected.Arguments = append(expected.Arguments[:i], expected.Arguments[i+1:]...)
			i--
		}
	}
	actual := a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{
		Conclusion: expected.Arguments[0].Conclusion,
	})
	assertArgumentSetsMatch(t, expected.Arguments, actual)
}

func TestGetWithOffsets(t *testing.T) {
	expected := parseGetAllResponse(t, acceptancetest.ReadFile(t, samplesPath+"get-all-response.json"))
	a := newApp(t, nil)
	a.SaveAllSuccessfully(t, expected.Arguments)
	fetched := a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{
		Count: 1,
	})
	assertArgumentSetsMatch(t, []arguments.Argument{expected.Arguments[0]}, fetched)
	fetched = a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{
		Count:  1,
		Offset: 1,
	})
	assertArgumentSetsMatch(t, []arguments.Argument{expected.Arguments[1]}, fetched)
}

func TestGetWithExclusions(t *testing.T) {
	expected := parseGetAllResponse(t, acceptancetest.ReadFile(t, samplesPath+"get-all-response.json"))
	a := newApp(t, nil)
	a.SaveAllSuccessfully(t, expected.Arguments)
	fetched := a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{
		Exclude: []int64{expected.Arguments[0].ID},
	})
	assertArgumentSetsMatch(t, expected.Arguments[1:], fetched)
	fetched = a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{
		Exclude: []int64{expected.Arguments[1].ID},
	})
	assertArgumentSetsMatch(t, append(append([]arguments.Argument{}, expected.Arguments[0]), expected.Arguments[2:]...), fetched)
}

func TestGetWithSearch(t *testing.T) {
	containing := []string{"bing", "words"}
	available := parseGetAllResponse(t, acceptancetest.ReadFile(t, samplesPath+"get-all-response.json"))
	expected := make([]arguments.Argument, 0)
	for i := 0; i < len(available.Arguments); i++ {
		if strings.Contains(available.Arguments[i].Conclusion, containing[0]) && strings.Contains(available.Arguments[i].Conclusion, containing[1]) {
			expected = append(expected, available.Arguments[i])
		}
	}
	a := newApp(t, nil)
	a.SaveAllSuccessfully(t, available.Arguments)
	fetched := a.FetchSomeSuccessfully(t, arguments.FetchSomeOptions{
		ConclusionContainsAll: containing,
	})
	assertArgumentSetsMatch(t, expected, fetched)
}

func parseGetAllResponse(t *testing.T, data []byte) argumentsHttp.GetAllResponse {
	var getAll argumentsHttp.GetAllResponse
	require.NoError(t, json.Unmarshal(data, &getAll))
	return getAll
}
