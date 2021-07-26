// Package synthesizers contains synthesizers (text to speech)
package synthesizers

// Synthesizer is a text to speech API interface
type Synthesizer interface {
	TextToSpeech(s string) ([]byte, error)
}
