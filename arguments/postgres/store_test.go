package postgres_test

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/wikisophia/api-arguments/arguments/argumentstest"

	"github.com/smotes/purse"
	"github.com/stretchr/testify/suite"
	argumentsInPostgres "github.com/wikisophia/api-arguments/arguments/postgres"
	"github.com/wikisophia/api-arguments/config"
	"github.com/wikisophia/api-arguments/postgres"
)

var hasDatabase = flag.Bool("database", false, "run database integration tests")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestArgumentStorageIntegration(t *testing.T) {
	// Only run tests which rely on the database if the database flag is present
	if !*hasDatabase {
		return
	}

	// These tests will be slow... so do as much as we can up front to save time
	sqlScripts, err := purse.New(filepath.Join("..", "..", "postgres", "scripts"))
	if err != nil {
		t.Fatalf("Failed to load file in purse: %v", err)
	}
	contents, ok := sqlScripts.Get("clear.sql")
	if !ok {
		t.Fatal("purse could not load storage/postgres/scritps/clear.sql")
	}

	db := postgres.NewDB(config.MustParse().Storage.Postgres)
	_, err = db.Query(contents)
	if err != nil {
		t.Fatalf("failed to execute queries in storage/postgres/scritps/clear.sql: %v", err)
	}

	store := argumentsInPostgres.NewStore(db)
	// Run all the same tests from the StoreTests suite.
	suite.Run(t, &DatabaseTests{
		StoreTests: &argumentstest.StoreTests{
			Store: store,
		},
		db:            db,
		queryContents: contents,
	})

	store.Close()
	db.Close()
}

type DatabaseTests struct {
	*argumentstest.StoreTests

	db            *sql.DB
	queryContents string
}

func (suite *DatabaseTests) SetupTest() {
	if suite.db != nil {
		_, err := suite.db.Query(suite.queryContents)
		if err != nil {
			fmt.Printf("failed to execute queries in storage/postgres/scritps/clear.sql: %v", err)
			os.Exit(1)
		}
	}
}
