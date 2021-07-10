// Package tokenizer contains tokenizers
package tokenizer

// Tokenizer takes a string and gets it's part of speech tokens
type Tokenizer interface {
	Setup() error
	Cleanup() error
	IsSetup() bool

	GetTokens(s string) ([]Token, error)
}

// Token is a part of speech token
type Token struct {
	POS        string `json:"pos"`
	EndIndex   uint   `json:"endIndex"`
	BeginIndex uint   `json:"beginIndex"`
	Morph      string `json:"morph"`
}

// NotSetupError is returned when a tokenizer function runs when it is not setup
type NotSetupError struct{}

func (e *NotSetupError) Error() string { return "tokenizer is not setup" }
