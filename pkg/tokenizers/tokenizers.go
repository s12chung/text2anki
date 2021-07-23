// Package tokenizers contains tokenizers
package tokenizers

import "github.com/s12chung/text2anki/pkg/lang"

// Tokenizer takes a string and gets it's part of speech tokens
type Tokenizer interface {
	Setup() error
	Cleanup() error
	IsSetup() bool

	GetTokens(s string) ([]Token, error)
}

// Token is a part of speech token
type Token struct {
	Text         string
	PartOfSpeech lang.PartOfSpeech
	StartIndex   uint
	EndIndex     uint
}

// NotSetupError is returned when a tokenizer function runs when it is not setup
type NotSetupError struct{}

func (e *NotSetupError) Error() string { return "tokenizer is not setup" }
