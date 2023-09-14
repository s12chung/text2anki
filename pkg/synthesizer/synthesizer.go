// Package synthesizer contains synthesizers (text to speech)
package synthesizer

import "context"

// Synthesizer is a text to speech API interface
type Synthesizer interface {
	SourceName() string
	TextToSpeech(ctx context.Context, text string) ([]byte, error)
}
