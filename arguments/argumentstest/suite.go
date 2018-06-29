package argumentstest

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api-arguments/arguments"
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
