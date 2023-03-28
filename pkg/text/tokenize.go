package text

import (
	"github.com/s12chung/text2anki/pkg/tokenizers"
)

// TokenizedText is the text grouped with it's tokens
type TokenizedText struct {
	Text
	Tokens []tokenizers.Token
}

// TokenizeTexts takes the texts and tokenizes them
func TokenizeTexts(tokenizer tokenizers.Tokenizer, texts []Text) (tokenizedTexts []TokenizedText, err error) {
	if err = tokenizer.Setup(); err != nil {
		return nil, err
	}
	defer func() {
		err2 := tokenizer.Cleanup()
		if err == nil {
			err = err2
		}
	}()

	tokenizedTexts = make([]TokenizedText, len(texts))
	for i, text := range texts {
		var tokens []tokenizers.Token
		tokens, err = tokenizer.Tokenize(text.Text)
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
