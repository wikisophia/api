package memory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api/server/arguments"
	"github.com/wikisophia/api/server/arguments/memory"
	"github.com/wikisophia/api/server/arguments/storetest"
)

// TestInMemoryStore makes sure that the inMemoryStore is consistent with the StoreTests suite.
// This helps verify:
//    1. The inMemoryStore, which is used throughout app tests to avoid a DB dependency.
//    2.The StoreTests suite, which is reused to test the real Postgres implementation.
func TestInMemoryStore(t *testing.T) {
	suite.Run(t, &storetest.StoreTests{
		StoreFactory: func() arguments.Store {
			return memory.NewMemoryStore()
		},
	})
}
