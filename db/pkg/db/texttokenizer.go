package db

import (
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
)

// TextTokenizer is used to generate TokenizedText
type TextTokenizer struct {
	Parser       text.Parser
	Tokenizer    tokenizers.Tokenizer
	CleanSpeaker bool
}

// TokenizedText is the text grouped with its tokens
type TokenizedText struct {
	text.Text
	Tokens []tokenizers.Token
}

// TokenizeTextsFromString converts a string to TokenizedText
func (t TextTokenizer) TokenizeTextsFromString(s string) ([]TokenizedText, error) {
	texts, err := t.Parser.TextsFromString(s)
	if err != nil {
		return nil, err
	}
	if t.CleanSpeaker {
		texts = text.CleanSpeaker(texts)
	}
	return t.TokenizeTexts(texts)
}

// TokenizeTexts takes the texts and tokenizes them
func (t TextTokenizer) TokenizeTexts(texts []text.Text) (tokenizedTexts []TokenizedText, err error) {
	if err = t.Tokenizer.Setup(); err != nil {
		return nil, err
	}
	defer func() {
		err2 := t.Tokenizer.Cleanup()
		if err == nil {
			err = err2
		}
	}()

	tokenizedTexts = make([]TokenizedText, len(texts))
	for i, text := range texts {
		var tokens []tokenizers.Token
		tokens, err = t.Tokenizer.Tokenize(text.Text)
		if err != nil {
			return nil, err
		}
		tokenizedTexts[i] = TokenizedText{
			Text:   text,
			Tokens: tokens,
		}
	}

	return tokenizedTexts, nil
}
