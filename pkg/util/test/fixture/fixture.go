// Package fixture contains helper functions for fixtures
package fixture

import (
	"crypto/sha256"
	"encoding/json"
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

	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

// TestDataDir returns the testdata dir
const TestDataDir = "testdata"

// UpdateFailMessage is a standard message to identify if the test is failed via. updating fixtures
//
// Wording may be used to clean up test message output
const UpdateFailMessage = updateFixturesEnv + "=true, fixtures are updated, turn off ENV var to run test"

// JoinTestData joins the elem path to the testdata dir
func JoinTestData(elem ...string) string {
	dirs := append([]string{TestDataDir}, elem...)
	return path.Join(dirs...)
}

// Read reads the fixture
func Read(t *testing.T, fixturePath string) []byte {
	require := require.New(t)
	expected, err := os.ReadFile(JoinTestData(fixturePath)) //nolint:gosec // for tests
	require.NoError(err)
	return []byte(strings.TrimSpace(string(expected)))
}

// Update updates fixture, used externally initial creation of test only
func Update(t *testing.T, fixturePath string, resultBytes []byte) {
	assert := assert.New(t)
	if !WillUpdate() {
		assert.Fail("fixtures.Update() is called without WillUpdate() == true")
	}

	err := os.MkdirAll(path.Join(TestDataDir, filepath.Dir(fixturePath)), os.ModePerm)
	assert.NoError(err) //nolint:testifylint // requires assert to keep making more fixtures

	err = os.WriteFile(JoinTestData(fixturePath), resultBytes, ioutil.OwnerRWGroupR)
	assert.NoError(err) //nolint:testifylint // requires assert to keep making more fixtures

	if WillUpdate() {
		assert.Fail(UpdateFailMessage) //nolint:testifylint // requires assert to keep making more fixtures
	}
}

// SafeUpdate calls Update(), only when WillUpdate() is true
func SafeUpdate(t *testing.T, fixturePath string, resultBytes []byte) {
	if WillUpdate() {
		Update(t, fixturePath, resultBytes)
	}
}

const envTrue = "true"
const updateFixturesEnv = "UPDATE_FIXTURES"

// WillUpdate returns true if the fixtures will be updated from ReadOrWrite
func WillUpdate() bool {
	return os.Getenv(updateFixturesEnv) == envTrue
}

// ReadOrUpdate reads the fixture or updates it if WillUpdate is true
func ReadOrUpdate(t *testing.T, fixturePath string, resultBytes []byte) []byte {
	if WillUpdate() {
		Update(t, fixturePath, resultBytes)
	}
	return Read(t, fixturePath)
}

// CompareReadOrUpdate calls ReadOrUpdate and compares the result against it
func CompareReadOrUpdate(t *testing.T, fixturePath string, resultBytes []byte) {
	require := require.New(t)
	expected := ReadOrUpdate(t, fixturePath, resultBytes)
	require.Equal(string(expected), strings.TrimSpace(string(resultBytes)))
}

// CompareReadOrUpdateJSON calls CompareReadOrUpdate, but adds .json to the fixturePath and calls JSON on the resultBytes
func CompareReadOrUpdateJSON(t *testing.T, fixturePath string, obj any) {
	CompareReadOrUpdate(t, fixturePath+".json", JSON(t, obj))
}

// CompareRead calls Read and compares the result against it
func CompareRead(t *testing.T, fixturePath string, resultBytes []byte) {
	require := require.New(t)
	expected := Read(t, fixturePath)
	require.Equal(string(expected), strings.TrimSpace(string(resultBytes)))
}

// CompareReadOrUpdateDir reads a fixture dir or updates it if WillUpdate is true
func CompareReadOrUpdateDir(t *testing.T, fixtureDir, resultDir string) {
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
		require.Fail(UpdateFailMessage)
	}
}

func compareOrUpdateDirName(t *testing.T, expected, result string) {
	require := require.New(t)

	if WillUpdate() {
		require.NoError(os.Mkdir(expected, ioutil.OwnerRWXGroupRX))
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
		require.NoError(os.WriteFile(expected, resultBytes, ioutil.OwnerRWGroupR)) //nolint:gosec // for tests
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

// JSON returns indented json for fixtures
func JSON(t *testing.T, v any) []byte {
	require := require.New(t)
	bytes, err := json.MarshalIndent(v, "", "  ")
	require.NoError(err)
	return bytes
}

// SHA2Map takes a directory path and generates a map between
// the filenames in the directory and their SHA2 hash.
//
// Often used with fixtures
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
