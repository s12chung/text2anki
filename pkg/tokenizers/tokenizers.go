// Package tokenizers contains tokenizers
package tokenizers

import (
	"fmt"

	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/tokenizers/server"
)

// Tokenizer takes a string and gets it's part of speech tokens
type Tokenizer interface {
	Setup() error
	Cleanup() error
	IsSetup() bool

	Tokenize(str string) ([]Token, error)
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

// ServerTokenizer is a wrapper around a Tokenizer implemented by an API server
type ServerTokenizer struct {
	name    string
	server  server.TokenizerServer
	started bool
}

// NewServerTokenizer returns a new ServerTokenizer
func NewServerTokenizer(name string, server server.TokenizerServer) ServerTokenizer {
	return ServerTokenizer{
		name:   name,
		server: server,
	}
}

// Setup setups up the JVM for Komoran to run
func (s *ServerTokenizer) Setup() error {
	if s.started {
		return fmt.Errorf("%v already started before, make a new instance", s.name)
	}
	s.started = true
	return s.server.Start()
}

// Cleanup cleans Komoran's resources
func (s *ServerTokenizer) Cleanup() error {
	return s.server.Stop()
}

// CleanupAndWait runs Cleanup() and waits until the server is stopped
func (s *ServerTokenizer) CleanupAndWait() error {
	return s.server.StopAndWait()
}

// IsSetup returns true if Komoran is ready to execute
func (s *ServerTokenizer) IsSetup() bool {
	return s.server.IsRunning()
}

// ServerTokenize returns the part of speech tokens of the given string from the server
func (s *ServerTokenizer) ServerTokenize(str string, resp any) error {
	if !s.IsSetup() {
		return &NotSetupError{}
	}

	return s.server.Tokenize(str, resp)
}
