// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

// SetupTempDB calls db.SetDB with a temp file
func SetupTempDB(testName string) error {
	filename := test.GenerateFilename(testName, ".sqlite3")
	if err := db.SetDB(path.Join(os.TempDir(), filename)); err != nil {
		return err
	}
	return db.Create(context.Background())
}

// SetupTempDBT calls SetupTempDB and checks errors
func SetupTempDBT(t *testing.T, testName string) {
	require := require.New(t)
	err := SetupTempDB(testName)
	require.NoError(err)
}

// Seed seeds the database with a small amount of data
func Seed(t *testing.T) {
	require := require.New(t)
	_, callerPath, _, ok := runtime.Caller(0)
	require.True(ok)

	bytes, err := os.ReadFile(path.Join(path.Dir(callerPath), fixture.TestDataDir, "Seed") + ".json")
	require.NoError(err)

	var terms []db.Term
	err = json.Unmarshal(bytes, &terms)
	require.NoError(err)

	queries := db.Qs()
	for _, term := range terms {
		_, err = queries.TermCreate(context.Background(), term.CreateParams())
		require.NoError(err)
	}
}

// SearchTerm is a search term used for tests
const SearchTerm = "마음"

// SearchConfig is the config used for test searching (so it stays constant)
var SearchConfig = db.TermsSearchConfig{
	PopLog:       20,
	PopWeight:    40,
	CommonWeight: 40,
	LenLog:       2,
}
