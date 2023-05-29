package db_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db/testdb"
)

type MustSetupAndSeed struct{}

func TestMain(m *testing.M) {
	testdb.MustSetupAndSeed(MustSetupAndSeed{})

	if err := textTokenizer.Setup(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	code := m.Run()
	if err := textTokenizer.Cleanup(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	os.Exit(code)
}

func testRecentTimestamps(t *testing.T, timestamps ...time.Time) {
	require := require.New(t)
	for _, timestamp := range timestamps {
		require.Greater(time.Now(), timestamp)
	}
}
