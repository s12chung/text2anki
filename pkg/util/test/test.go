// Package test contains test helpers
package test

import (
	"os"
	"testing"
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
