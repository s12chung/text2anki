// Package stringclean have a set of string cleaning utilities
package stringclean

import (
	"strings"
)

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
