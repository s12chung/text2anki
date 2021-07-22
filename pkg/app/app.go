// Package app contains app specific functions
package app

import (
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/dictionary"
)

// NewNoteFromTerm returns a note given the Term and index
func NewNoteFromTerm(term dictionary.Term, translationIndex uint) anki.Note {
	translation := term.Translations[translationIndex]
	return anki.Note{
		Text:         term.Text,
		PartOfSpeech: term.PartOfSpeech,
		Translation:  translation.Text,

		CommonLevel:      term.CommonLevel,
		Explanation:      translation.Explanation,
		DictionarySource: term.DictionarySource,
	}
}
