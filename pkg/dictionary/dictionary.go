// Package dictionary contains dictonary related functions
package dictionary

import "github.com/s12chung/text2anki/pkg/lang"

// Dicionary is a interface for a dictionary to search with
type Dicionary interface {
	// Search returns the dictionary terms for the given query in popularity order
	Search(q string) ([]Term, error)
}

// Term is a word or phrase
type Term struct {
	Text         string
	Variants     []string
	PartOfSpeech lang.PartOfSpeech
	CommonLevel  lang.CommonLevel
	Translations []Translation

	DictionarySource string
}

// Translation is a translation of a word or phrase
type Translation struct {
	Text        string
	Explanation string
}
