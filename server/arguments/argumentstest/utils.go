package argumentstest

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments"
)

// AssertArgumentsMatch fails if the arguments aren't equivalent.
// Similar to assert.Equal, but ignores Premise order.
func AssertArgumentsMatch(t *testing.T, expected arguments.Argument, actual arguments.Argument) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Conclusion, actual.Conclusion)
	assert.ElementsMatch(t, expected.Premises, actual.Premises)
}

// ParseSample parses the JSON from the file at unixPath and returns
// it as an Argument.
func ParseSample(t *testing.T, unixPath string) arguments.Argument {
	return ParseJSON(t, ReadFile(t, unixPath))
}

// ReadFile reads the JSON data from unixPath.
func ReadFile(t *testing.T, unixPath string) []byte {
	fileBytes, err := ioutil.ReadFile(filepath.FromSlash(unixPath))
	assert.NoError(t, err)
	return fileBytes
}

// ParseJSON unmarshals JSON data as an Argument.
func ParseJSON(t *testing.T, data []byte) arguments.Argument {
	var argument arguments.Argument
	assert.NoError(t, json.Unmarshal(data, &argument))
	return argument
}
