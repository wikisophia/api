package memory

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/wikisophia/api-arguments/server/arguments"
)

// NewStore returns an in-memory implementation of a Store.
//
// This is used when testing other parts of the app so that those tests don't
// need to rely on a database.
//
// TODO #11: This should be threadsafe. It's not a huge deal yet because this
// is used for tests & development... but might cause some false positives.
func NewStore() arguments.Store {
	// Populate the arguments value with a "dummy" arg, since versions start at 1.
	// The implementation is just a bit simpler if we start the real data at index 1 too.
	return &inMemoryStore{
		nextID:    0,
		arguments: make([][]arguments.Argument, 1),
	}
}

type inMemoryStore struct {
	nextID    int64
	arguments [][]arguments.Argument
}

func (s *inMemoryStore) Delete(ctx context.Context, id int64) error {
	if id < int64(len(s.arguments)) {
		s.arguments[id] = nil
	}
	return nil
}

func (s *inMemoryStore) FetchVersion(ctx context.Context, id int64, version int16) (arguments.Argument, error) {
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

func (s *inMemoryStore) FetchLive(ctx context.Context, id int64) (arguments.Argument, error) {
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

func (s *inMemoryStore) FetchAll(ctx context.Context, conclusion string) ([]arguments.ArgumentWithID, error) {
	args := make([]arguments.ArgumentWithID, 0, 20)
	for i := 1; i < len(s.arguments); i++ {
		if s.arguments[i][0].Conclusion == conclusion {
			args = append(args, arguments.ArgumentWithID{
				Argument: s.arguments[i][len(s.arguments[i])-1],
				ID:       int64(i),
			})
		}
	}
	return args, nil
}

func (s *inMemoryStore) Save(ctx context.Context, argument arguments.Argument) (id int64, err error) {
	s.arguments = append(s.arguments, []arguments.Argument{
		argument, // Add this twice because the 0th index will be ignored by Fetches
		argument,
	})
	return int64(len(s.arguments) - 1), nil
}

func (s *inMemoryStore) UpdatePremises(ctx context.Context, argumentID int64, premises []string) (version int16, err error) {
	if !s.argumentExists(argumentID) {
		return -1, &arguments.NotFoundError{
			Message: fmt.Sprintf("argument with id %d does not exist", argumentID),
		}
	}
	conclusion := s.arguments[argumentID][0].Conclusion
	s.arguments[argumentID] = append(s.arguments[argumentID], arguments.Argument{
		Conclusion: conclusion,
		Premises:   premises,
	})
	return int16(len(s.arguments[argumentID]) - 1), nil
}

func (s *inMemoryStore) argumentExists(id int64) bool {
	return int64(len(s.arguments)) > id && s.arguments[id] != nil
}

func (s *inMemoryStore) Close() error {
	return nil
}

func findMax(data map[int16]arguments.Argument) (max int16) {
	for key := range data {
		if max < key {
			max = key
		}
	}
	return
}
