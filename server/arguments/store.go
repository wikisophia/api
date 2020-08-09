package arguments

import (
	"context"
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
