// Package stringutil have a set of string utilities
package stringutil

import (
	"strings"
)

// SplitClean splits the string with the given separator and trims the space for each element
func SplitClean(s, sep string) []string {
	slc := []string{}
	if strings.TrimSpace(s) != "" {
		slc = strings.Split(s, sep)
	}
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return slc
}
