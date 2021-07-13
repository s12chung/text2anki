// Package fixture contains helper functions for fixtures
package fixture

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// WillUpdate returns true if the fixtures will be updated from ReadOrWrite
func WillUpdate() bool {
	return os.Getenv("UPDATE_FIXTURES") == "true"
}

// ReadOrUpdate reads the fixture or updates it if WillUpdate is true
func ReadOrUpdate(t *testing.T, fixtureFilename string, resultBytes []byte) string {
	require := require.New(t)
	if WillUpdate() {
		err := ioutil.WriteFile(fixtureFilename, resultBytes, 0600)
		require.Nil(err)
		require.FailNow("UPDATE_FIXTURES=true, fixtures are updated, turn off ENV var to run test")
	}
	//nolint:gosec // for tests
	expected, err := ioutil.ReadFile(fixtureFilename)
	require.Nil(err)
	return strings.TrimSpace(string(expected))
}

// CompareReadOrUpdate calls ReadOrUpdate and compares the result against it
func CompareReadOrUpdate(t *testing.T, fixtureFilename string, resultBytes []byte) {
	require := require.New(t)
	expected := ReadOrUpdate(t, fixtureFilename, resultBytes)
	require.Equal(string(resultBytes), expected)
}
