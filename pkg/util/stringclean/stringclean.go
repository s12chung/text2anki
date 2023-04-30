// Package stringclean have a set of string cleaning utilities
package stringclean

import (
	"regexp"
	"strings"
)

var speakerRegex = regexp.MustCompile(`\A[^:\d]{0,25}:`)

// Speaker removes the "speaker name:" string from s
func Speaker(s string) string {
	s = speakerRegex.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

// Split splits the string with the given separator and trims the space for each element
func Split(s, sep string) []string {
	slc := []string{}
	if strings.TrimSpace(s) != "" {
		slc = strings.Split(s, sep)
	}
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return slc
}
