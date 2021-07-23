// Package text contains functions to separate text into source and translation
package text

import (
	"strings"
)

// Text represents a text line given from the source
type Text struct {
	Text        string
	Translation string
}

// TextsFromString returns an array of Texts from the given string
func TextsFromString(s string) []Text {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	texts := make([]Text, 0, len(lines))
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		texts = append(texts, Text{Text: l})
	}
	return texts
}
