package csv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestFile(t *testing.T) {
	require := require.New(t)
	file, err := os.CreateTemp("", "text2ankiTestFile-*.csv")
	require.NoError(err)
	defer func() { require.NoError(os.Remove(file.Name())) }()

	err = File(file.Name(), [][]string{
		{"test", "me"},
		{"1", "2"},
	})
	require.NoError(err)

	bytes, err := os.ReadFile(file.Name())
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestFile.csv", bytes)
}
