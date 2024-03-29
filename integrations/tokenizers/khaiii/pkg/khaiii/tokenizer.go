package khaiii

import (
	"log/slog"

	"github.com/s12chung/text2anki/pkg/util/logg"
)

// Tokenizer is a wrapper around the Khaiii
type Tokenizer struct {
	kahiii *Khaiii
	log    *slog.Logger
}

// NewTokenizer returns a new Tokenizer
func NewTokenizer(k *Khaiii, log *slog.Logger) *Tokenizer {
	return &Tokenizer{
		kahiii: k,
		log:    log,
	}
}

// Cleanup cleans up the Kahiii instance
func (k *Tokenizer) Cleanup() {
	if err := k.kahiii.Close(); err != nil {
		k.log.Error("Kahiii.Cleanup()", logg.Err(err))
	}
}

// TokenizeResponse is the response given by Tokenize
type TokenizeResponse struct {
	Words []Word `json:"words"`
}

// Tokenize returns the tokenized words of the given string
func (k *Tokenizer) Tokenize(str string) (any, error) {
	words, err := k.kahiii.Analyze(str)
	if err != nil {
		return nil, err
	}
	return &TokenizeResponse{Words: words}, nil
}
