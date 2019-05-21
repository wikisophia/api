package argumentstest

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/endpoints"
)

// StoreTests is a testing suite which makes sure that a Store obeys
// the interface contract
type StoreTests struct {
	suite.Suite
	StoreFactory func() endpoints.Store
}

// TestSaveIsLive makes sure that an argument is "live" immediately after being saved.
func (suite *StoreTests) TestSaveIsLive() {
	original := ParseSample(suite.T(), "../../samples/save-request.json")
	store := suite.StoreFactory()
	id, err := store.Save(context.Background(), original)
	if !assert.NoError(suite.T(), err) {
		return
	}
	original.ID = id
	original.Version = 1
	fetched, err := store.FetchLive(context.Background(), id)
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Equal(suite.T(), original, fetched)
}

// TestUpdatedIsLive makes sure that a newly updated argument uses the latest premises.
func (suite *StoreTests) TestUpdatedIsLive() {
	original := ParseSample(suite.T(), "../../samples/save-request.json")
	updated := ParseSample(suite.T(), "../../samples/update-request.json")

	store := suite.StoreFactory()
	id := suite.saveWithUpdates(store, original, updated)
	if id == -1 {
		return
	}
	updated.ID = id
	updated.Version = 2
	fetched, err := store.FetchLive(context.Background(), id)
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Equal(suite.T(), updated, fetched)
}

// TestUpdateUnknownReturnsError makes sure that we can't update arguments which don't exist.
func (suite *StoreTests) TestUpdateUnknownReturnsError() {
	store := suite.StoreFactory()
	unknown := ParseSample(suite.T(), "../../samples/save-request.json")
	unknown.ID = 1
	_, err := store.Update(context.Background(), unknown)
	if !assert.Error(suite.T(), err) {
		return
	}

	if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.UpdatePremises() should return a NotFoundError on arguments which don't exist.")
	}
}

// TestOriginalIsAvailable makes sure that old versions of updated arguments can still be fetched.
func (suite *StoreTests) TestOriginalIsAvailable() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../../samples/save-request.json")
	updated := ParseSample(suite.T(), "../../samples/update-request.json")

	id := suite.saveWithUpdates(store, original, updated)
	if id == -1 {
		return
	}
	original.ID = id
	original.Version = 1
	fetched, err := store.FetchVersion(context.Background(), id, 1)
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Equal(suite.T(), original, fetched)
}

// TestDeletedUnknownReturnsNotFound makes sure the backend returns a NotFoundError
// if asked to delete an unknown entry.
func (suite *StoreTests) TestDeletedUnknownReturnsNotFound() {
	store := suite.StoreFactory()
	err := store.Delete(context.Background(), 1)
	if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.Delete() should return a NotFoundError for unknown IDs.")
	}
}

// TestDeletedIsUnavailable makes sure the backend doesn't return arguments that have been deleted.
func (suite *StoreTests) TestDeletedIsUnavailable() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../../samples/save-request.json")
	updated := ParseSample(suite.T(), "../../samples/update-request.json")

	id := suite.saveWithUpdates(store, original, updated)
	if id == -1 {
		return
	}
	if !assert.NoError(suite.T(), store.Delete(context.Background(), id)) {
		return
	}
	if _, err := store.FetchVersion(context.Background(), id, 1); !assert.Error(suite.T(), err) {
		return
	}
	if _, err := store.FetchLive(context.Background(), id); !assert.Error(suite.T(), err) {
		return
	} else if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store should return a NotFoundError on deleted arguments.")
	}
}

// TestFetchUnknownReturnsError makes sure the backend returns errors when asked for an unknown ID.
func (suite *StoreTests) TestFetchUnknownReturnsError() {
	store := suite.StoreFactory()
	if _, err := store.FetchLive(context.Background(), 1); !assert.Error(suite.T(), err) {
		return
	} else if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.FetchLive should return a NotFoundError on unknown arguments.")
	}
	if _, err := store.FetchVersion(context.Background(), 1, 1); !assert.Error(suite.T(), err) {
		return
	} else if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.FetchVersion should return a NotFoundError on unknown arguments.")
	}
}

// TestBasicFetchAll makes sure the Store returns all the arguments for a conclusion.
func (suite *StoreTests) TestBasicFetchAll() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../../samples/save-request.json")
	updated := ParseSample(suite.T(), "../../samples/update-request.json")

	original.ID = suite.saveWithUpdates(store, original)
	otherArg := arguments.Argument{
		Conclusion: original.Conclusion,
		Premises:   updated.Premises,
	}
	otherArg.ID = suite.saveWithUpdates(store, otherArg)

	suite.saveWithUpdates(store, arguments.Argument{
		Conclusion: "some other conclusion",
		Premises:   []string{"premise1", "premise2"},
	})

	allArgs, err := store.FetchAll(context.Background(), original.Conclusion)

	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 2) {
		return
	}
	original.Version = 1
	otherArg.Version = 1

	// Fixes #1: Arguments might be returned in any order
	fetchedFirst := allArgs[0]
	fetchedSecond := allArgs[1]
	if fetchedFirst.ID != original.ID {
		tmp := fetchedFirst
		fetchedFirst = fetchedSecond
		fetchedSecond = tmp
	}

	assert.Equal(suite.T(), original, fetchedFirst)
	assert.Equal(suite.T(), otherArg, fetchedSecond)
}

// TestVersionedFetchAll makes sure the Store returns the argument's live version only.
func (suite *StoreTests) TestVersionedFetchAll() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../../samples/save-request.json")
	updated := ParseSample(suite.T(), "../../samples/update-request.json")
	updated.Conclusion = original.Conclusion

	id := suite.saveWithUpdates(store, original, updated)
	allArgs, err := store.FetchAll(context.Background(), original.Conclusion)
	updated.ID = id
	updated.Version = 2
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 1) {
		return
	}
	assert.Equal(suite.T(), updated, allArgs[0])
}

func (suite *StoreTests) saveWithUpdates(store endpoints.Store, arg arguments.Argument, updates ...arguments.Argument) int64 {
	id, err := store.Save(context.Background(), arg)
	if !assert.NoError(suite.T(), err) {
		return -1
	}
	for i := 0; i < len(updates); i++ {
		update := updates[i]
		update.ID = id
		_, err = store.Update(context.Background(), update)
		if !assert.NoError(suite.T(), err) {
			return -1
		}
	}

	return id
}
