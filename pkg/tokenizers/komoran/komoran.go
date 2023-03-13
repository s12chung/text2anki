// Package komoran is an interface to the Komoran Korean tokenizer
package komoran

import (
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/server"
	"github.com/s12chung/text2anki/pkg/tokenizers/server/java"
)

// New returns a Komoran Korean tokenizer
func New() tokenizers.Tokenizer {
	return &Komoran{
		server: java.NewJarServer(jarPath, 9999, 64),
	}
}

// Komoran is a Korean Tokenizer in java
type Komoran struct {
	server server.Server
}

var jarPath = "tokenizers/dist/komoran/tokenizer-komoran.jar"

// Setup setups up the JVM for Komoran to run
func (k *Komoran) Setup() error {
	return k.server.Start()
}

// Cleanup cleans Komoran's resources
func (k *Komoran) Cleanup() error {
	return k.server.Stop()
}

// IsSetup returns true if Komoran is ready to execute
func (k *Komoran) IsSetup() bool {
	return k.server.IsRunning()
}

// Tokenize returns the part of speech tokens of the given string
func (k *Komoran) Tokenize(str string) ([]tokenizers.Token, error) {
	if !k.IsSetup() {
		return nil, &tokenizers.NotSetupError{}
	}

	resp := &response{}
	err := k.server.Tokenize(str, resp)
	if err != nil {
		return nil, err
	}
	return resp.toTokenizerTokens(), nil
}
