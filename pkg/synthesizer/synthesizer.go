// Package synthesizer contains synthesizers (text to speech)
package synthesizer

import "context"

// Synthesizer is a text to speech API interface
type Synthesizer interface {
	TextToSpeech(ctx context.Context, s string) ([]byte, error)
	SourceName() string
}
