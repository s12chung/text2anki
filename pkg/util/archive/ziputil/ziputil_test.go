package ziputil

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestZipDir(t *testing.T) {
	require := require.New(t)
	testName := "TestZipDir"

	b, err := ZipDir(fixture.JoinTestData(testName, "testdir"))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, filepath.Join(testName, "result.zip"), b)
}
