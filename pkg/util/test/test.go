// Package test contains test helpers
package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

// IsCI returns true if the test envirnment is under CI
func IsCI() bool {
	return os.Getenv("CI") == "true"
}

// CISkip skips the tests if IsCI == "true"
func CISkip(t *testing.T, msg string) {
	if IsCI() {
		t.Skip(msg)
	}
}

var timeNow = time.Now

// GenerateFilename returns a non-colliding filename
func GenerateFilename(name, ext string) string {
	if ext[0] != '.' {
		ext = "." + ext
	}
	return GenerateName(name) + ext
}

// GenerateName returns a non-colliding name
func GenerateName(name string) string {
	return fmt.Sprintf("text2anki-%v-%v", name, timeNow().Format(time.StampNano))
}

// MkdirAll is a simple wrapper around os.MkdirAll
func MkdirAll(t *testing.T, path string) {
	require.NoError(t, os.MkdirAll(path, ioutil.OwnerRWXGroupRX))
}
