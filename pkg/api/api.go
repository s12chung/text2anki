// Package api contains the app configured structures
package api

import (
	"os"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizers/komoran"
)

// TextTokenizer is the TextTokenizer used in the app
var TextTokenizer = db.TextTokenizer{
	Parser:       DefaultParser(),
	Tokenizer:    DefaultTokenizer(),
	CleanSpeaker: true,
}

// DefaultParser returns the default Parser used in the app
func DefaultParser() text.Parser {
	return text.NewParser(text.Korean, text.English)
}

// DefaultTokenizer returns the default Tokenizer used in the app
func DefaultTokenizer() tokenizers.Tokenizer {
	switch os.Getenv("TOKENIZER") {
	case "komoran":
		return komoran.New()
	default:
		return khaiii.New()
	}
}
