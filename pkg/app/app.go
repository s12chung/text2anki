// Package app contains app specific functions
package app

import (
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
)

// TokenizedText is the text grouped with it's tokens
type TokenizedText struct {
	text.Text
	Tokens []tokenizers.Token
}

// TokenizeTexts takes the texts and tokenizes them
func TokenizeTexts(texts []text.Text) (tokenizedTexts []TokenizedText, err error) {
	// tokenizer := komoran.New()
	// if err = tokenizer.Setup(); err != nil {
	// 	return nil, err
	// }
	// defer func() {
	// 	err2 := tokenizer.Cleanup()
	// 	if err == nil {
	// 		err = err2
	// 	}
	// }()

	// tokenizedTexts = make([]TokenizedText, len(texts))
	// for i, text := range texts {
	// 	var tokens []tokenizers.Token
	// 	tokens, err = tokenizer.GetTokens(text.Text)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	tokenizedTexts[i] = TokenizedText{
	// 		Text:   text,
	// 		Tokens: tokens,
	// 	}
	// }

	return tokenizedTexts, nil
}

// NewNoteFromTerm returns a Note given the Term and index
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
