// Package app contains app specific functions
package app

import (
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/komoran"
)

// TokenizedText is the text grouped with it's tokens
type TokenizedText struct {
	text.Text
	Tokens []tokenizers.Token
}

// TokenizeTexts takes the texts and tokenizes them
func TokenizeTexts(texts []text.Text) ([]TokenizedText, error) {
	tokenizer := komoran.NewKomoran()
	var err2 error
	if err := tokenizer.Setup(); err != nil {
		return nil, err
	}
	defer func() {
		err2 = tokenizer.Cleanup()
	}()

	tokenizedTexts := make([]TokenizedText, len(texts))
	for i, text := range texts {
		tokens, err := tokenizer.GetTokens(text.Text)
		if err != nil {
			return nil, err
		}
		tokenizedTexts[i] = TokenizedText{
			Text:   text,
			Tokens: tokens,
		}
	}

	return tokenizedTexts, err2
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

// NewNoteFromText returns a Note given the given Text
func NewNoteFromText(text text.Text) anki.Note {
	return anki.Note{
		Text:             text.Text,
		PartOfSpeech:     lang.PartOfSpeechUnknown,
		Translation:      text.Translation,
		DictionarySource: "Text2Anki Imported Text",
	}
}
