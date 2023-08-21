package tokenizer

import (
	"strings"

	"github.com/s12chung/text2anki/pkg/lang"
)

// SplitTokenizer is a simple tokenizer that splits the string with a space, used for testing or mocks
type SplitTokenizer struct {
	setup bool
}

// NewSplitTokenizer returns a new SplitTokenizer
func NewSplitTokenizer() Tokenizer {
	return &SplitTokenizer{}
}

// Setup makes setup = true
func (s *SplitTokenizer) Setup() error {
	s.setup = true
	return nil
}

// Cleanup makes setup = false
func (s *SplitTokenizer) Cleanup() error {
	s.setup = false
	return nil
}

// IsSetup  returns setup
func (s *SplitTokenizer) IsSetup() bool {
	return s.setup
}

// Tokenize returns tokens set as nouns, split by a space
func (s *SplitTokenizer) Tokenize(str string) ([]Token, error) {
	split := strings.Split(str, " ")
	tokens := make([]Token, len(split))
	var index uint = 0
	for i, s := range split {
		length := uint(len(s))
		tokens[i] = Token{
			Text:         s,
			PartOfSpeech: lang.PartOfSpeechNoun,
			StartIndex:   index,
			Length:       length,
		}
		index += length + 1
	}
	return tokens, nil
}
