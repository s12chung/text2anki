// Package test contains test helpers
package test

import (
	"encoding/json"
	"fmt"
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

// MustJSON marshals v into indented JSON, panics if fails
func MustJSON(v any) []byte {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("MustJSON: error marshaling: %v", v))
	}
	return bytes
}
