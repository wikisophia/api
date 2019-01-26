package argumentstest

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api-arguments/server/arguments"
)

// StoreTests is a testing suite which makes sure that a Store obeys
// the interface contract
type StoreTests struct {
	suite.Suite
	Store arguments.Store
}

var originalArguments = arguments.Argument{
	Conclusion: "Socrates is mortal",
	Premises: []string{
		"Socrates is a human",
		"All men are mortal",
	},
}

var updatedPremises = []string{
	"Socrates is a man",
	"All men are mortal",
}

// TestSaveIsLive makes sure that an argument is "live" immediately after being saved.
func (suite *StoreTests) TestSaveIsLive() {
	id, err := suite.Store.Save(context.Background(), originalArguments)
	if !assert.NoError(suite.T(), err) {
		return
	}

	fetched, err := suite.Store.FetchLive(context.Background(), id)
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Equal(suite.T(), originalArguments.Conclusion, fetched.Conclusion)
	assert.ElementsMatch(suite.T(), originalArguments.Premises, fetched.Premises)
}

// TestUpdatedIsLive makes sure that a newly updated argument uses the latest premises.
func (suite *StoreTests) TestUpdatedIsLive() {
	id := suite.saveWithUpdates(originalArguments, updatedPremises)
	if id == -1 {
		return
	}

	fetched, err := suite.Store.FetchLive(context.Background(), id)
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Equal(suite.T(), originalArguments.Conclusion, fetched.Conclusion)
	assert.ElementsMatch(suite.T(), updatedPremises, fetched.Premises)
}

// TestUpdateUnknownReturnsError makes sure that we can't update arguments which don't exist.
func (suite *StoreTests) TestUpdateUnknownReturnsError() {
	_, err := suite.Store.UpdatePremises(context.Background(), 1, []string{"Socrates is a man", "All men are mortal"})
	if !assert.Error(suite.T(), err) {
		return
	}

	if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.UpdatePremises() should return a NotFoundError on arguments which don't exist.")
	}
}

// TestOriginalIsAvailable makes sure that old versions of updated arguments can still be fetched.
func (suite *StoreTests) TestOriginalIsAvailable() {
	id := suite.saveWithUpdates(originalArguments, updatedPremises)
	if id == -1 {
		return
	}

	fetched, err := suite.Store.FetchVersion(context.Background(), id, 1)
	if !assert.NoError(suite.T(), err) {
		return
	}
	assert.Equal(suite.T(), originalArguments.Conclusion, fetched.Conclusion)
	assert.ElementsMatch(suite.T(), originalArguments.Premises, fetched.Premises)
}

// TestDeletedIsUnavailable makes sure the backend doesn't return arguments that have been deleted.
func (suite *StoreTests) TestDeletedIsUnavailable() {
	id := suite.saveWithUpdates(originalArguments, updatedPremises)
	if id == -1 {
		return
	}
	if !assert.NoError(suite.T(), suite.Store.Delete(context.Background(), id)) {
		return
	}
	if _, err := suite.Store.FetchVersion(context.Background(), id, 1); !assert.Error(suite.T(), err) {
		return
	}
	if _, err := suite.Store.FetchLive(context.Background(), id); !assert.Error(suite.T(), err) {
		return
	} else if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store should return a NotFoundError on deleted arguments.")
	}
}

// TestFetchUnknownReturnsError makes sure the backend returns errors when asked for an unknown ID.
func (suite *StoreTests) TestFetchUnknownReturnsError() {
	if _, err := suite.Store.FetchLive(context.Background(), 1); !assert.Error(suite.T(), err) {
		return
	} else if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.FetchLive should return a NotFoundError on unknown arguments.")
	}
	if _, err := suite.Store.FetchVersion(context.Background(), 1, 1); !assert.Error(suite.T(), err) {
		return
	} else if _, ok := err.(*arguments.NotFoundError); !ok {
		suite.T().Error("Store.FetchVersion should return a NotFoundError on unknown arguments.")
	}
}

// TestBasicFetchAll makes sure the Store returns all the arguments for a conclusion.
func (suite *StoreTests) TestBasicFetchAll() {
	arg1ID := suite.saveWithUpdates(originalArguments)
	otherArg := arguments.Argument{
		Conclusion: originalArguments.Conclusion,
		Premises:   updatedPremises,
	}
	suite.saveWithUpdates(arguments.Argument{
		Conclusion: "some other conclusion",
		Premises:   []string{"premise1", "premise2"},
	})
	arg2ID := suite.saveWithUpdates(otherArg)

	allArgs, err := suite.Store.FetchAll(context.Background(), originalArguments.Conclusion)

	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 2) {
		return
	}

	// Fixes #1: Arguments might be returned in any order
	fetchedFirst := allArgs[0]
	fetchedSecond := allArgs[1]
	if fetchedFirst.ID != arg1ID {
		tmp := fetchedFirst
		fetchedFirst = fetchedSecond
		fetchedSecond = tmp
	}

	assert.Equal(suite.T(), originalArguments.Conclusion, fetchedFirst.Conclusion)
	assert.ElementsMatch(suite.T(), originalArguments.Premises, fetchedFirst.Premises)
	assert.Equal(suite.T(), arg1ID, fetchedFirst.ID)
	assert.Equal(suite.T(), originalArguments.Conclusion, fetchedSecond.Conclusion)
	assert.ElementsMatch(suite.T(), updatedPremises, fetchedSecond.Premises)
	assert.Equal(suite.T(), arg2ID, fetchedSecond.ID)
}

// TestVersionedFetchAll makes sure the Store returns the argument's live version only.
func (suite *StoreTests) TestVersionedFetchAll() {
	arg1ID := suite.saveWithUpdates(originalArguments, updatedPremises)
	allArgs, err := suite.Store.FetchAll(context.Background(), originalArguments.Conclusion)
	if !assert.NoError(suite.T(), err) {
		return
	}
	if !assert.Len(suite.T(), allArgs, 1) {
		return
	}
	assert.Equal(suite.T(), originalArguments.Conclusion, allArgs[0].Conclusion)
	assert.ElementsMatch(suite.T(), updatedPremises, allArgs[0].Premises)
	assert.Equal(suite.T(), arg1ID, allArgs[0].ID)
}

func (suite *StoreTests) saveWithUpdates(arg arguments.Argument, premiseUpdates ...[]string) int64 {
	id, err := suite.Store.Save(context.Background(), arg)
	if !assert.NoError(suite.T(), err) {
		return -1
	}

	for i := 0; i < len(premiseUpdates); i++ {
		_, err = suite.Store.UpdatePremises(context.Background(), id, premiseUpdates[i])
		if !assert.NoError(suite.T(), err) {
			return -1
		}
	}

	return id
}
