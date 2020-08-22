package arguments

import (
	"context"
	"fmt"
)

// Store combines all the functions needed to read & write Arguments
// into a single interface.
type Store interface {
	Deleter
	GetSome
	GetVersioned
	GetLive
	Saver
	Updater
	Close() error
}

// Deleter can delete arguments by ID.
type Deleter interface {
	// Delete deletes an argument (and all its versions) from the site.
	// If the argument didn't exist, the error will be a NotFoundError.
	Delete(ctx context.Context, id int64) error
}

// GetSome can fetch lists of arguments at once.
type GetSome interface {
	// FetchSome finds the arguments which match the options.
	// If none exist, error will be nil and the slice empty.
	FetchSome(ctx context.Context, options FetchSomeOptions) ([]Argument, error)
}

// GetVersioned returns a specific version of an argument.
type GetVersioned interface {
	// FetchVersion should return a particular version of an argument.
	// If the the argument didn't exist, the error should be an arguments.NotFoundError.
	FetchVersion(ctx context.Context, id int64, version int) (Argument, error)
}

// GetLive can fetch the live version of an argument.
type GetLive interface {
	// FetchLive should return the latest "active" version of an argument.
	// If no argument with this ID exists, the error should be an arguments.NotFoundError.
	FetchLive(ctx context.Context, id int64) (Argument, error)
}

// Saver can save arguments.
type Saver interface {
	// Save stores an argument and returns that argument's ID.
	// The ID on the input argument will be ignored.
	Save(ctx context.Context, argument Argument) (id int64, err error)
}

// Updater can update existing arguments.
type Updater interface {
	// Update makes a new version of the argument. It returns the new argument's version.
	// If no argument with this ID exists, the returned error is an arguments.NotFoundError.
	Update(ctx context.Context, argument Argument) (version int, err error)
}

// FetchSomeOptions has some ways to limit what gets returned when fetching all the arguments.
type FetchSomeOptions struct {
	// Conclusion only finds arguments which support a given conclusion
	Conclusion string
	// ConclusionContainsAll limits returned arguments to ones with conclusions that
	// contain all the words in this array.
	ConclusionContainsAll []string
	// Count limits the number of fetched arguments.
	Count int
	// Exclude prevents arguments which have any of these IDs from being returned
	Exclude []int64
	// Offset changes which arguments start being returned.
	//
	// An offset of 0 will return arguments starting with the first one.
	// An offset of 1 will skip the first argument, and return arguments starting with the second.
	//
	// When combined with Count, this can be used to paginate the results.
	Offset int
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
