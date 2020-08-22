package memory

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/wikisophia/api/server/arguments"
)

// NewMemoryStore returns an in-memory implementation of a Store.
//
// This is used when testing other parts of the app so that those tests don't
// need to rely on a database.
//
// TODO #11: This should be threadsafe. It's not a huge deal yet because this
// is used for tests & development... but might cause some false positives.
func NewMemoryStore() *InMemoryStore {
	// Populate the arguments value with a "dummy" arg, since versions start at 1.
	// The implementation is just a bit simpler if we start the real data at index 1 too.
	return &InMemoryStore{
		nextID:    0,
		arguments: make([][]arguments.Argument, 1),
	}
}

// InMemoryStore saves arguments in program memory.
// This is mainly intended for testing and easier dev environment setups.
type InMemoryStore struct {
	nextID    int64
	arguments [][]arguments.Argument
}

// Delete deletes an argument (and all its versions) from the site.
// If the argument didn't exist, the error will be a NotFoundError.
func (s *InMemoryStore) Delete(ctx context.Context, id int64) error {
	if id > 0 && id < int64(len(s.arguments)) {
		s.arguments[id] = nil
		return nil
	}
	return &arguments.NotFoundError{
		Message: fmt.Sprintf("argument with id %d does not exist", id),
	}
}

// FetchVersion should return a particular version of an argument.
// If the the argument didn't exist, the error should be an NotFoundError.
func (s *InMemoryStore) FetchVersion(ctx context.Context, id int64, version int) (arguments.Argument, error) {
	if !s.argumentExists(id) {
		return arguments.Argument{}, &arguments.NotFoundError{
			Message: fmt.Sprintf("argument with id %d does not exist", id),
		}
	}
	versions := s.arguments[id]
	if len(versions) <= int(version) {
		return arguments.Argument{}, &arguments.NotFoundError{
			Message: fmt.Sprintf("version %d of argument %d does not exist", version, id),
		}
	}
	return versions[version], nil
}

// FetchLive should return the latest "active" version of an argument.
// If no argument with this ID exists, the error should be an NotFoundError.
func (s *InMemoryStore) FetchLive(ctx context.Context, id int64) (arguments.Argument, error) {
	if !s.argumentExists(id) {
		return arguments.Argument{}, &arguments.NotFoundError{
			Message: fmt.Sprintf("argument with id %d does not exist", id),
		}
	}
	versions := s.arguments[id]
	if len(versions) == 1 {
		return arguments.Argument{}, errors.New("versions was empty... this is a bug in the InMemoryStore test code")
	}
	return versions[len(versions)-1], nil
}

// FetchSome returns all the "live" arguments matching the given options.
// If none exist, error will be nil and the slice empty.
func (s *InMemoryStore) FetchSome(ctx context.Context, options arguments.FetchSomeOptions) ([]arguments.Argument, error) {
	args := make([]arguments.Argument, 0, 20)
	numSkipped := 0
	for i := 1; i < len(s.arguments); i++ {
		if containsInt64(options.Exclude, int64(i)) {
			continue
		}
		if options.Conclusion != "" && options.Conclusion != s.arguments[i][len(s.arguments[i])-1].Conclusion {
			continue
		}
		if !containsAll(s.arguments[i][len(s.arguments[i])-1].Conclusion, options.ConclusionContainsAll) {
			continue
		}

		if numSkipped >= options.Offset {
			args = append(args, s.arguments[i][len(s.arguments[i])-1])
			if len(args) == options.Count {
				return args, nil
			}
		} else {
			numSkipped++
		}
	}
	return args, nil
}

func containsAll(text string, elements []string) bool {
	for _, element := range elements {
		if !strings.Contains(text, element) {
			return false
		}
	}
	return true
}

// Save stores an argument and returns that argument's ID.
// The ID on the input argument will be ignored.
func (s *InMemoryStore) Save(ctx context.Context, argument arguments.Argument) (id int64, err error) {
	argument.ID = int64(len(s.arguments))
	argument.Version = 1
	s.arguments = append(s.arguments, []arguments.Argument{
		argument, // Add this twice because the 0th index will be ignored by Fetches
		argument,
	})
	return argument.ID, nil
}

// Update makes a new version of the argument. It returns the new argument's version.
// If no argument with this ID exists, the returned error is an NotFoundError.
func (s *InMemoryStore) Update(ctx context.Context, argument arguments.Argument) (version int, err error) {
	if !s.argumentExists(argument.ID) {
		return -1, &arguments.NotFoundError{
			Message: fmt.Sprintf("argument with id %d does not exist", argument.ID),
		}
	}
	argument.Version = len(s.arguments[argument.ID])
	s.arguments[argument.ID] = append(s.arguments[argument.ID], argument)
	return argument.Version, nil
}

// Close frees all the resources the InMemoryStore is using (needed to implement the interface)
func (s *InMemoryStore) Close() error {
	return nil
}

func (s *InMemoryStore) argumentExists(id int64) bool {
	return int64(len(s.arguments)) > id && s.arguments[id] != nil
}

func containsInt64(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
