// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
)

// SetupTempDB calls db.SetDB with a temp file
func SetupTempDB(ctx context.Context, testName string) error {
	filename := fmt.Sprintf("text2anki-%v-%v.sqlite3", testName, time.Now().Unix())
	if err := db.SetDB(path.Join(os.TempDir(), filename)); err != nil {
		return err
	}
	return db.Create(ctx)
}

// SetupTempDBT calls SetupTempDB and checks errors
func SetupTempDBT(ctx context.Context, t *testing.T, testName string) {
	require := require.New(t)
	err := SetupTempDB(ctx, testName)
	require.NoError(err)
}
