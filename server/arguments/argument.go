package arguments

import (
	"errors"
	"fmt"
)

// Argument is the core data type for the API.
type Argument struct {
	ID         int64    `json:"id"`
	Version    int      `json:"version"`
	Conclusion string   `json:"conclusion"`
	Premises   []string `json:"premises"`
}

// NotFoundError will be returned by Store.Fetch() calls when the cause of the returned error is
// that the argument simply doesn't exist.
type NotFoundError struct {
	Message string
	Args    []interface{}
}

// FetchSomeOptions has some ways to limit what gets returned when fetching all the arguments.
type FetchSomeOptions struct {
	// Conclusion only finds arguments which support a given conclusion
	Conclusion string
	// Count limits the number of fetched objects.
	Count int
	// Offset changes which arguments start being returned.
	//
	// An offset of 0 will return arguments starting with the first one.
	// An offset of 1 will skip the first argument, and return arguments starting with the second.
	//
	// When combined with Count, this can be used to paginate the results.
	Offset int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf(e.Message, e.Args...)
}

// ValidateArgument returns nil if the given argument has all the required pieces
// (e.g. non-nil conclusion + premises), and an error if the given argument is malformed.
func ValidateArgument(argument Argument) error {
	if argument.Conclusion == "" {
		return errors.New("arguments must have a conclusion")
	}
	if len(argument.Premises) < 2 {
		return errors.New("arguments must have at least 2 premises")
	}
	for i, premise := range argument.Premises {
		if premise == "" {
			return fmt.Errorf("argument premise[%d] is empty, but must not be", i)
		}
	}

	return nil
}
