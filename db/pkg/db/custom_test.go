package db

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/storage/localstore/localstoretest"
	"github.com/s12chung/text2anki/pkg/util/test"
)

func init() {
	api, err := localstoretest.NewAPI("db.custom_test")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	SetDBStorage(storage.NewDBStorage(api, nil))
}

func dBPath(testName string) string {
	return path.Join(os.TempDir(), test.GenerateFilename(testName, ".sqlite3"))
}

func TestSetDB(t *testing.T) {
	oldDB := database
	defer func() {
		database = oldDB
	}()

	require := require.New(t)
	err := SetDB(dBPath("TestSetDB"))
	require.NoError(err)
}

func TestCreate(t *testing.T) {
	oldDB := database
	defer func() {
		database = oldDB
	}()

	require := require.New(t)
	err := SetDB(dBPath("TestCreate"))
	require.NoError(err)
	err = Qs().Create(context.Background())
	require.NoError(err)
}
