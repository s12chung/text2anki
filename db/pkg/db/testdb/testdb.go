// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

// SetupTempDB calls db.SetDB with a temp file
func SetupTempDB(testName string) error {
	filename := fmt.Sprintf("text2anki-%v-%v.sqlite3", testName, time.Now().Unix())
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
