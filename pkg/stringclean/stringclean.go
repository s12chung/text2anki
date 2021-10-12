// Package stringclean have a set of string cleaning utilities
package stringclean

import (
	"regexp"
	"strings"
)

var speakerRegex = regexp.MustCompile("\\A[^:]{0,25}:")

// Speaker removes the "speaker name:" string from s
func Speaker(s string) string {
	s = speakerRegex.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}
