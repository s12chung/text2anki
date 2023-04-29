package db

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func dBPath(testName string) string {
	filename := fmt.Sprintf("text2anki-%v-%v.sqlite3", testName, time.Now().Unix())
	return path.Join(os.TempDir(), filename)
}

func TestSetDB(t *testing.T) {
	require := require.New(t)
	err := SetDB(dBPath("TestSetDB"))
	require.NoError(err)
}

func TestCreate(t *testing.T) {
	require := require.New(t)
	err := SetDB(dBPath("TestCreate"))
	require.NoError(err)
	err = Create(context.Background())
	require.NoError(err)
}
