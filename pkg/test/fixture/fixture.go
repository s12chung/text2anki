// Package fixture contains helper functions for fixtures
package fixture

import (
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
	//nolint:gosec // for tests
	expected, err := os.ReadFile(JoinTestData(fixtureFilename))
	require.Nil(err)
	return []byte(strings.TrimSpace(string(expected)))
}

// Update updates fixture, used externally initial creation of test only
func Update(t *testing.T, fixtureFilename string, resultBytes []byte) {
	assert := assert.New(t)
	if !WillUpdate() {
		assert.Fail("fixtures.Update() is called without WillUpdate() == true")
	}

	err := os.MkdirAll(TestDataDir, os.ModePerm)
	assert.Nil(err)

	err = os.WriteFile(JoinTestData(fixtureFilename), resultBytes, 0600)
	assert.Nil(err)

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
		require.Nil(os.RemoveAll(fixtureDir))
	}

	err := filepath.WalkDir(resultDir, func(result string, d fs.DirEntry, err error) error {
		require.Nil(err)
		rel, err := filepath.Rel(resultDir, result)
		require.Nil(err)

		expected := path.Join(fixtureDir, rel)
		if d.IsDir() {
			compareOrUpdateDirName(t, expected, result)
		} else {
			compareOrUpdateFile(t, expected, result)
		}
		return nil
	})
	require.Nil(err)
	if WillUpdate() {
		require.Fail(fmt.Sprintf("%v=true, fixtures are updated, turn off ENV var to run test", updateFixturesEnv))
	}
}

func compareOrUpdateDirName(t *testing.T, expected, result string) {
	require := require.New(t)

	if WillUpdate() {
		require.Nil(os.Mkdir(expected, 0750))
		return
	}

	stat, err := os.Stat(expected)
	require.Nil(err)
	require.True(stat.IsDir(), fmt.Sprintf("result, %v, is not matching %v", result, expected))
}

func compareOrUpdateFile(t *testing.T, expected, result string) {
	require := require.New(t)

	//nolint:gosec // for tests
	resultBytes, err := os.ReadFile(result)
	require.Nil(err)

	if WillUpdate() {
		//nolint:gosec // for tests
		require.Nil(os.WriteFile(expected, resultBytes, 0600))
		return
	}

	//nolint:gosec // for tests
	expectedBytes, err := os.ReadFile(expected)
	require.Nil(err)

	if utf8.Valid(expectedBytes) && utf8.Valid(resultBytes) {
		require.Equal(string(expectedBytes), string(resultBytes))
	} else {
		require.Equal(expectedBytes, resultBytes)
	}
}
