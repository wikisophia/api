package argumentstest

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api/server/arguments"
)

// StoreTests is a testing suite which makes sure that a Store obeys
// the interface contract
type StoreTests struct {
	suite.Suite
	StoreFactory func() arguments.Store
}

// TestSaveIsLive makes sure that an argument is "live" immediately after being saved.
func (suite *StoreTests) TestSaveIsLive() {
	original := ParseSample(suite.T(), "../samples/save-request.json")
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
	original := ParseSample(suite.T(), "../samples/save-request.json")
	updated := ParseSample(suite.T(), "../samples/update-request.json")

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
	unknown := ParseSample(suite.T(), "../samples/save-request.json")
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
	original := ParseSample(suite.T(), "../samples/save-request.json")
	updated := ParseSample(suite.T(), "../samples/update-request.json")

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
	original := ParseSample(suite.T(), "../samples/save-request.json")
	updated := ParseSample(suite.T(), "../samples/update-request.json")

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

// TestFetchByConclusion makes sure the Store returns all the arguments for a conclusion.
func (suite *StoreTests) TestFetchByConclusion() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../samples/save-request.json")
	updated := ParseSample(suite.T(), "../samples/update-request.json")

	original.ID = suite.saveWithUpdates(store, original)
	original.Version = 1

	otherArg := arguments.Argument{
		Conclusion: original.Conclusion,
		Premises:   updated.Premises,
	}
	otherArg.ID = suite.saveWithUpdates(store, otherArg)
	otherArg.Version = 1

	suite.saveWithUpdates(store, arguments.Argument{
		Conclusion: "some other conclusion",
		Premises:   []string{"premise1", "premise2"},
	})

	allArgs, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Conclusion: original.Conclusion,
	})

	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 2) {
		return
	}

	assert.Equal(suite.T(), original, allArgs[0])
	assert.Equal(suite.T(), otherArg, allArgs[1])
}

// TestFetchWithConclusionSearch makes sure the Store limits what it returns
// based on which words the user expects the conclusion to have.
func (suite *StoreTests) TestFetchWithConclusionSearch() {
	store := suite.StoreFactory()
	sample := ParseSample(suite.T(), "../samples/save-request.json")

	first := suite.saveCopyWithConclusion(store, sample, "best of times")
	second := suite.saveCopyWithConclusion(store, sample, "worst of times")
	suite.saveCopyWithConclusion(store, sample, "unrelated conclusion")

	found, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		ConclusionContainsAll: []string{"times"},
	})
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), found, 2) {
		return
	}
	assert.Equal(suite.T(), first, found[0])
	assert.Equal(suite.T(), second, found[1])

	found, err = store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		ConclusionContainsAll: []string{"best", "times"},
	})
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), found, 1) {
		return
	}
	assert.Equal(suite.T(), first, found[0])
}

func (suite *StoreTests) saveCopyWithConclusion(store arguments.Store, template arguments.Argument, conclusion string) arguments.Argument {
	copy := template
	copy.Conclusion = conclusion
	copy.ID = suite.saveWithUpdates(store, copy)
	copy.Version = 1
	return copy
}

// TestVersionedFetchAll makes sure the Store returns the argument's live version only.
func (suite *StoreTests) TestVersionedFetchAll() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../samples/save-request.json")
	updated := ParseSample(suite.T(), "../samples/update-request.json")
	updated.Conclusion = original.Conclusion

	id := suite.saveWithUpdates(store, original, updated)
	allArgs, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Conclusion: original.Conclusion,
	})
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

// TestFetchAllChangedConclusion makes sure the store finds live versions
// which have a different conclusion from when they started.
func (suite *StoreTests) TestFetchAllChangedConclusion() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../samples/save-request.json")
	updated := ParseSample(suite.T(), "../samples/update-request.json")

	id := suite.saveWithUpdates(store, original, updated)
	allArgs, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Conclusion: updated.Conclusion,
	})
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

// TestFetchOne makes sure the Store limits how many objects it returns properly.
func (suite *StoreTests) TestFetchOne() {
	store := suite.StoreFactory()
	original := ParseSample(suite.T(), "../samples/save-request.json")

	suite.saveWithUpdates(store, original)
	suite.saveWithUpdates(store, original)

	allArgs, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Count: 1,
	})
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Len(suite.T(), allArgs, 1)
}

// TestFetchWithOffset makes sure the Store skips elements properly when given an offset.
func (suite *StoreTests) TestFetchWithOffset() {
	store := suite.StoreFactory()

	first := ParseSample(suite.T(), "../samples/save-request.json")
	suite.saveWithUpdates(store, first)
	second := arguments.Argument{
		Conclusion: "some second conclusion",
		Premises:   first.Premises,
	}
	suite.saveWithUpdates(store, second)
	third := arguments.Argument{
		Conclusion: "some third conclusion",
		Premises:   first.Premises,
	}
	suite.saveWithUpdates(store, third)

	fetchAndAssert := func(offset int, conclusion string) {
		allArgs, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
			Count:  1,
			Offset: offset,
		})
		if !assert.NoError(suite.T(), err) {
			return
		}
		assert.Len(suite.T(), allArgs, 1)
		assert.Equal(suite.T(), conclusion, allArgs[0].Conclusion)
	}

	fetchAndAssert(0, first.Conclusion)
	fetchAndAssert(1, second.Conclusion)
	fetchAndAssert(2, third.Conclusion)
}

// TestFetchWithExclusions makes sure the Store excludes arguments properly.
func (suite *StoreTests) TestFetchWithExclusions() {
	store := suite.StoreFactory()
	arg := ParseSample(suite.T(), "../samples/save-request.json")

	id1 := suite.saveWithUpdates(store, arg)
	id2 := suite.saveWithUpdates(store, arg)

	allArgs, err := store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Exclude: []int64{id1},
	})
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 1) {
		return
	}
	assert.Equal(suite.T(), id2, allArgs[0].ID)

	allArgs, err = store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Exclude: []int64{id2},
	})
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 1) {
		return
	}
	assert.Equal(suite.T(), id1, allArgs[0].ID)

	id3 := suite.saveWithUpdates(store, arg)
	allArgs, err = store.FetchSome(context.Background(), arguments.FetchSomeOptions{
		Exclude: []int64{id2},
	})
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 2) {
		return
	}
	assert.Equal(suite.T(), id1, allArgs[0].ID)
	assert.Equal(suite.T(), id3, allArgs[1].ID)
}

func (suite *StoreTests) saveWithUpdates(store arguments.Store, arg arguments.Argument, updates ...arguments.Argument) int64 {
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
