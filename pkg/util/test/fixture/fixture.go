// Package fixture contains helper functions for fixtures
package fixture

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDataDir returns the testdata dir
const TestDataDir = "testdata"

// JoinTestData joins the elem path to the testdata dir
func JoinTestData(elem ...string) string {
	dirs := append([]string{TestDataDir}, elem...)
	return path.Join(dirs...)
}

// Read reads the fixture
func Read(t *testing.T, fixtureFilename string) []byte {
	require := require.New(t)
	expected, err := os.ReadFile(JoinTestData(fixtureFilename)) //nolint:gosec // for tests
	require.NoError(err)
	return []byte(strings.TrimSpace(string(expected)))
}

// Update updates fixture, used externally initial creation of test only
func Update(t *testing.T, fixtureFilename string, resultBytes []byte) {
	assert := assert.New(t)
	if !WillUpdate() {
		assert.Fail("fixtures.Update() is called without WillUpdate() == true")
	}

	err := os.MkdirAll(TestDataDir, os.ModePerm)
	assert.NoError(err)

	err = os.WriteFile(JoinTestData(fixtureFilename), resultBytes, 0600)
	assert.NoError(err)

	if WillUpdate() {
		assert.Fail(fmt.Sprintf("%v=true, fixtures are updated, turn off ENV var to run test", updateFixturesEnv))
	}
}

// SafeUpdate calls Update(), only when WillUpdate() is true
func SafeUpdate(t *testing.T, fixtureFilename string, resultBytes []byte) {
	if WillUpdate() {
		Update(t, fixtureFilename, resultBytes)
	}
}

const envTrue = "true"
const updateFixturesEnv = "UPDATE_FIXTURES"

// WillUpdate returns true if the fixtures will be updated from ReadOrWrite
func WillUpdate() bool {
	return os.Getenv(updateFixturesEnv) == envTrue
}

// ReadOrUpdate reads the fixture or updates it if WillUpdate is true
func ReadOrUpdate(t *testing.T, fixtureFilename string, resultBytes []byte) []byte {
	if WillUpdate() {
		Update(t, fixtureFilename, resultBytes)
	}
	return Read(t, fixtureFilename)
}

// CompareReadOrUpdate calls ReadOrUpdate and compares the result against it
func CompareReadOrUpdate(t *testing.T, fixtureFilename string, resultBytes []byte) {
	require := require.New(t)
	expected := ReadOrUpdate(t, fixtureFilename, resultBytes)
	require.Equal(string(expected), strings.TrimSpace(string(resultBytes)))
}

// CompareRead calls Read and compares the result against it
func CompareRead(t *testing.T, fixtureFilename string, resultBytes []byte) {
	require := require.New(t)
	expected := Read(t, fixtureFilename)
	require.Equal(string(expected), strings.TrimSpace(string(resultBytes)))
}

// CompareOrUpdateDir reads a fixture dir or updates it if WillUpdate is true
func CompareOrUpdateDir(t *testing.T, fixtureDir, resultDir string) {
	require := require.New(t)

	fixtureDir = JoinTestData(fixtureDir)
	if WillUpdate() {
		require.NoError(os.RemoveAll(fixtureDir))
	}

	err := filepath.WalkDir(resultDir, func(result string, d fs.DirEntry, err error) error {
		require.NoError(err)
		rel, err := filepath.Rel(resultDir, result)
		require.NoError(err)

		expected := path.Join(fixtureDir, rel)
		if d.IsDir() {
			compareOrUpdateDirName(t, expected, result)
		} else {
			compareOrUpdateFile(t, expected, result)
		}
		return nil
	})
	require.NoError(err)
	if WillUpdate() {
		require.Fail(fmt.Sprintf("%v=true, fixtures are updated, turn off ENV var to run test", updateFixturesEnv))
	}
}

func compareOrUpdateDirName(t *testing.T, expected, result string) {
	require := require.New(t)

	if WillUpdate() {
		require.NoError(os.Mkdir(expected, 0750))
		return
	}

	stat, err := os.Stat(expected)
	require.NoError(err)
	require.True(stat.IsDir(), fmt.Sprintf("result, %v, is not matching %v", result, expected))
}

func compareOrUpdateFile(t *testing.T, expected, result string) {
	require := require.New(t)

	resultBytes, err := os.ReadFile(result) //nolint:gosec // for tests
	require.NoError(err)

	if WillUpdate() {
		require.NoError(os.WriteFile(expected, resultBytes, 0600)) //nolint:gosec // for tests
		return
	}

	expectedBytes, err := os.ReadFile(expected) //nolint:gosec // for tests
	require.NoError(err)

	if utf8.Valid(expectedBytes) && utf8.Valid(resultBytes) {
		require.Equal(string(expectedBytes), string(resultBytes))
	} else {
		require.Equal(expectedBytes, resultBytes)
	}
}

// SHA2Map takes a directory path and generates a map between
// the filenames in the directory and their SHA2 hash.
//
// Often used with fixstures
func SHA2Map(dir string) (map[string]string, error) {
	fileMap := make(map[string]string)
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(p) //nolint:gosec // for tests
		if err != nil {
			return err
		}
		fileMap[d.Name()] = fmt.Sprintf("%x", sha256.Sum256(content))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileMap, nil
}
