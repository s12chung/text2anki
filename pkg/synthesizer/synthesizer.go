// Package synthesizer contains synthesizers (text to speech)
package synthesizer

// Synthesizer is a text to speech API interface
type Synthesizer interface {
	TextToSpeech(s string) ([]byte, error)
	SourceName() string
}
