package acceptancetest

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api/server/arguments"
)

// ParseSample parses the JSON from the file at unixPath and returns
// it as an Argument.
func ParseSample(t *testing.T, unixPath string) arguments.Argument {
	return ParseJSON(t, ReadFile(t, unixPath))
}

// ReadFile reads the JSON data from unixPath.
func ReadFile(t *testing.T, unixPath string) []byte {
	fileBytes, err := ioutil.ReadFile(filepath.FromSlash(unixPath))
	require.NoError(t, err)
	return fileBytes
}

// ParseJSON unmarshals JSON data as an Argument.
func ParseJSON(t *testing.T, data []byte) arguments.Argument {
	var argument arguments.Argument
	require.NoError(t, json.Unmarshal(data, &argument))
	return argument
}
