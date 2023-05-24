// Package test contains test helpers
package test

import (
	"fmt"
	"os"
	"testing"
	"time"
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
	return fmt.Sprintf("text2anki-%v-%v%v", name, timeNow().Format(time.StampMilli), ext)
}
