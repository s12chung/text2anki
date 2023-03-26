// Package komoran is an interface to the Komoran Korean tokenizer
package komoran

import (
	"fmt"
	"os"
	"time"

	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/server"
	"github.com/s12chung/text2anki/pkg/tokenizers/server/java"
)

// New returns a Komoran Korean tokenizer
func New() tokenizers.Tokenizer {
	return new()
}

const stopWarningDuration = 15 * time.Second

func new() *Komoran {
	return &Komoran{
		server: java.NewJarServer(jarPath, 9999, 64, stopWarningDuration),
	}
}

// Komoran is a Korean Tokenizer in java
type Komoran struct {
	server  server.TokenizerServer
	started bool
}

var jarPath string

func init() {
	jarPath = os.Getenv("KOMORAN_JAR_PATH")
}

// Setup setups up the JVM for Komoran to run
func (k *Komoran) Setup() error {
	if k.started {
		return fmt.Errorf("Komoran already started before, make a new instance")
	}
	k.started = true
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

type response struct {
	Tokens []token `json:"tokens"`
}

type token struct {
	POS        string `json:"pos"`
	EndIndex   uint   `json:"endIndex"`
	BeginIndex uint   `json:"beginIndex"`
	Morph      string `json:"morph"`
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
	return toTokenizerTokens(resp)
}
