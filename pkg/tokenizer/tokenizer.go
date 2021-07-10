// Package tokenizer contains tokenizers
package tokenizer

// Tokenizer takes a string and gets it's grammar tokens
type Tokenizer interface {
	Setup() error
	Cleanup() error
	IsSetup() bool

	GetTokens() (string, error)
}

// NotSetupError is returned when a tokenizer function runs when it is not setup
type NotSetupError struct{}

func (e *NotSetupError) Error() string { return "tokenizer is not setup" }
