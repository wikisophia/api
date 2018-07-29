package arguments

import (
	"context"
	"errors"
	"fmt"
)

type Argument struct {
	Conclusion string   `json:"conclusion"`
	Premises   []string `json:"premises"`
}

type ArgumentFromAll struct {
	Argument
	ID int64 `json:"id"`
}

// A Store manages Arguments inside the database.
type Store interface {
	// Delete deletes an argument (and all its versions) from the site.
	// This should only return an error if the deletion failed.
	// It should return nil if asked to delete an argument which doesn't exist.
	Delete(ctx context.Context, id int64) error
	// FetchAll pulls all the available arguments for a conclusion.
	// If none exist, error will be nil and the array empty.
	FetchAll(ctx context.Context, conclusion string) ([]ArgumentFromAll, error)
	// FetchVersion pulls a particular version of an argument from the database.
	// If the query completed successfully, but the argument didn't exist, the error
	// will be a NotFoundError.
	FetchVersion(ctx context.Context, id int64, version int16) (Argument, error)
	// FetchLatest pulls the latest version of an argument from the database.
	// If the query completed successfully, but the argument didn't exist, the error
	// will be a NotFoundError.
	FetchLive(ctx context.Context, id int64) (Argument, error)
	// Save stores an argument in the database, and returns that argument's ID.
	Save(ctx context.Context, argument Argument) (id int64, err error)
	// Update makes a new version of the argument with the given ID. It returns the argument version.
	// If the query completed successfully, but the original argument didn't exist, the error
	// will be a NotFoundError.
	UpdatePremises(ctx context.Context, argumentID int64, premises []string) (version int16, err error)
	// Close closes the store, freeing up its resources. Once called, the other functions
	// on the Store will fail.
	Close() error
}

type NotFoundError struct {
	Message string
	Args    []interface{}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf(e.Message, e.Args...)
}

func ValidateArgument(argument Argument) error {
	if argument.Conclusion == "" {
		return errors.New("arguments must have a conclusion")
	}
	if err := ValidatePremises(argument.Premises); err != nil {
		return err
	}
	return nil
}

func ValidatePremises(premises []string) error {
	if len(premises) < 2 {
		return errors.New("arguments must have at least 2 premises")
	}
	for i, premise := range premises {
		if premise == "" {
			return fmt.Errorf("argument premise[%d] is empty, but must not be", i)
		}
	}
	return nil
}
