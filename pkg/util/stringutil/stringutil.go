// Package stringutil have a set of string utilities
package stringutil

import (
	"strings"
	"unicode"
	"unicode/utf8"
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

// FirstUnbrokenSubstring gets the first unbroken phrase, under the given maxLength and TrimSpaces it
func FirstUnbrokenSubstring(s string, maxLength int) string {
	index := FirstUnbrokenIndex(s)
	if index != -1 {
		s = s[:index]
	}
	if index <= maxLength {
		return strings.TrimSpace(s)
	}
	index = strings.LastIndex(s, " ")
	if index == -1 {
		return strings.TrimSpace(s[:maxLength])
	}
	return strings.TrimSpace(s[:index])
}

var brokenCharacters = map[rune]bool{
	'!': true,
	',': true,
	'.': true,
	':': true,
	';': true,
	'?': true,
	'-': true,
	'–': true,
	'—': true,
}

// FirstUnbrokenIndex returns the first index that represents an unbroken phrase
func FirstUnbrokenIndex(s string) int {
	for i, char := range s {
		if !(unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsPunct(char) || unicode.IsSymbol(char) ||
			char == '_' || char == ' ') || brokenCharacters[char] {
			return i
		}
	}
	return -1
}

// TrimBytes trims the utf-8 string to the maxBytes
func TrimBytes(s string, maxBytes int) string {
	byteCount := 0
	for i := range s {
		_, size := utf8.DecodeRuneInString(s[i:])
		byteCount += size
		if byteCount > maxBytes {
			return s[:i]
		}
	}
	return s
}
