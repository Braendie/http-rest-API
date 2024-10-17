package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

// TestMain sets up the environment for running tests.
// It retrieves the "DATABASE_URL" environment variable. If it's not set, a default value for the test database is used.
// After setting up the database URL, it executes the tests by calling m.Run().
// The program then exits with the result of the test run.
func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=restapi_test user=braendie password=1234 sslmode=disable"
	}

	os.Exit(m.Run())
}
