package arguments

import (
	"context"
	"errors"
	"fmt"
)

// Argument is the core data type for the API.
type Argument struct {
	ID         int64    `json:"id"`
	Conclusion string   `json:"conclusion"`
	Premises   []string `json:"premises"`
}

// A Store manages Arguments inside the database.
type Store interface {
	// Delete deletes an argument (and all its versions) from the site.
	// If the argument didn't exist, the error will be a NotFoundError.
	Delete(ctx context.Context, id int64) error
	// FetchAll pulls all the available arguments for a conclusion.
	// If none exist, error will be nil and the array empty.
	FetchAll(ctx context.Context, conclusion string) ([]Argument, error)
	// FetchVersion pulls a particular version of an argument from the database.
	// If the query completed successfully, but the argument didn't exist, the error
	// will be a NotFoundError.
	FetchVersion(ctx context.Context, id int64, version int16) (Argument, error)
	// FetchLatest pulls the latest version of an argument from the database.
	// If the query completed successfully, but the argument didn't exist, the error
	// will be a NotFoundError.
	FetchLive(ctx context.Context, id int64) (Argument, error)
	// Save stores an argument in the database, and returns that argument's ID.
	// The ID on the incoming argument will be ignored.
	Save(ctx context.Context, argument Argument) (id int64, err error)
	// Update makes a new version of the argument. It returns the argument version.
	// If the returned error is a NotFoundError, then no argument with this ID exists,
	// and Save() should be used instead.
	Update(ctx context.Context, argument Argument) (version int16, err error)
	// Close closes the store, freeing up its resources.
	// Once called, the other functions on the Store will fail.
	Close() error
}

// NotFoundError will be returned by Store.Fetch() calls when the cause of the returned error is
// that the argument simply doesn't exist.
type NotFoundError struct {
	Message string
	Args    []interface{}
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
