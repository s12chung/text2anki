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

	zipBytes, err := ZipDir(fixture.JoinTestData(testName))
	require.NoError(err)

	paths := []string{
		"blah.txt",
		"innerdir/",
		"innerdir/waka.txt",
		"ok.txt",
	}
	require.NoError(CompareContents(zipBytes, paths, func(path string, contents []byte) {
		fixture.CompareReadOrUpdate(t, filepath.Join(testName, path), contents)
	}))
}
