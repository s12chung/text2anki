package khaiii

import (
	"fmt"
)

// Tokenizer is a wrapper around the Khaiii
type Tokenizer struct {
	kahiii *Khaiii
}

// NewTokenizer returns a new Tokenizer
func NewTokenizer(k *Khaiii) *Tokenizer {
	return &Tokenizer{
		kahiii: k,
	}
}

// Cleanup cleans up the Kahiii instance
func (k *Tokenizer) Cleanup() {
	if err := k.kahiii.Close(); err != nil {
		fmt.Println(err)
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
