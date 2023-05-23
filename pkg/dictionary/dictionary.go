// Package dictionary contains dictonary related functions
package dictionary

import "github.com/s12chung/text2anki/pkg/lang"

// Dictionary is an interface for a dictionary to search with
type Dictionary interface {
	// Search returns the dictionary terms for the given query in popularity order
	Search(q string, pos lang.PartOfSpeech) ([]Term, error)
}

// Term is a word or phrase
//
//nolint:musttag // Used for Temp UI only for YAML
type Term struct {
	Text         string            `json:"text,omitempty"`
	Variants     []string          `json:"variants,omitempty"`
	PartOfSpeech lang.PartOfSpeech `json:"part_of_speech,omitempty"`
	CommonLevel  lang.CommonLevel  `json:"common_level,omitempty"`
	Translations []Translation     `json:"translations,omitempty"`

	DictionarySource string `json:"dictionary_source,omitempty"`
}

// Translation is a translation of a word or phrase
type Translation struct {
	Text        string `json:"text,omitempty"`
	Explanation string `json:"explanation,omitempty"`
}
