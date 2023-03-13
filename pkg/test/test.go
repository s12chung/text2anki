// Package test contains test helpers
package test

import (
	"os"
)

// IsCI returns true if the test envirnment is under CI
func IsCI() bool {
	return os.Getenv("CI") == "true"
}
