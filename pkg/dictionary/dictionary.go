// Package dictionary contains dictonary related functions
package dictionary

import (
	"context"

	"github.com/s12chung/text2anki/pkg/lang"
)

// Dictionary is an interface for a dictionary to search with
type Dictionary interface {
	// Search returns the dictionary terms for the given query in popularity order
	Search(ctx context.Context, q string, pos lang.PartOfSpeech) ([]Term, error)
}

// Term is a word or phrase
type Term struct {
	ID           int64             `json:"id"`
	Text         string            `json:"text"`
	Variants     []string          `json:"variants"`
	PartOfSpeech lang.PartOfSpeech `json:"part_of_speech"`
	CommonLevel  lang.CommonLevel  `json:"common_level"`
	Translations []Translation     `json:"translations"`

	DictionarySource string `json:"dictionary_source"`
}

// StaticCopy returns a copy without fields that variate
func (t Term) StaticCopy() Term {
	c := t
	c.ID = 0
	return c
}

// Translation is a translation of a word or phrase
type Translation struct {
	Text        string `json:"text"`
	Explanation string `json:"explanation"`
}
