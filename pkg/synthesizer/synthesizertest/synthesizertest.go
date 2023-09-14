// Package synthesizertest provides a testing synthesizer
package synthesizertest

import "context"

// Synthesizer is a test Synthesizer
type Synthesizer struct{}

// SourceName returns the name of the Synthesizer
func (s Synthesizer) SourceName() string { return "synthesizertest" }

// TextToSpeech returns the byte representation of the text
func (s Synthesizer) TextToSpeech(_ context.Context, text string) ([]byte, error) {
	return []byte(text), nil
}
