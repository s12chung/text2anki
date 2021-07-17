// Package komoran is an interface to the Komoran Korean tokenizer
package komoran

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/pkg/java"
	"github.com/s12chung/text2anki/pkg/tokenizer"
	"tekao.net/jnigi"
)

// NewKomoran returns a Komoran Korean tokenizer
func NewKomoran() tokenizer.Tokenizer {
	return &Komoran{}
}

// Komoran is a Korean Tokenizer in java
type Komoran struct {
	javaInstance java.Instance
}

// Setup setups up the JVM for Komoran to run
func (k *Komoran) Setup() error {
	classPathArray, err := jarPaths()
	if err != nil {
		return err
	}
	return k.javaInstance.Setup(strings.Join(append(classPathArray, ""), ":"))
}

// Cleanup cleans Komoran's resources
func (k *Komoran) Cleanup() error {
	return k.javaInstance.Cleanup()
}

// IsSetup returns true if Komoran is ready to execute
func (k *Komoran) IsSetup() bool {
	return k.javaInstance.IsSetup()
}

// GetTokens returns the part of speech tokens of the given string
func (k *Komoran) GetTokens(s string) ([]tokenizer.Token, error) {
	if !k.javaInstance.IsSetup() {
		return nil, &tokenizer.NotSetupError{}
	}
	tokensJSON, err := k.callGetTokens(s)
	if err != nil {
		return nil, err
	}
	tokens := &komoranTokens{}
	if err := json.Unmarshal([]byte(tokensJSON), tokens); err != nil {
		return nil, err
	}
	return tokens.TokenList, nil
}

type komoranTokens struct {
	TokenList []tokenizer.Token `json:"tokenList"`
}

func (k *Komoran) callGetTokens(s string) (string, error) {
	jS, err := k.javaInstance.Env.NewObject("java/lang/String", []byte(s))
	if err != nil {
		return "", err
	}

	jTokensString, err := k.javaInstance.Env.CallStaticMethod("text2anki/tokenizer/komoran/Tokenizer",
		"getTokens",
		jnigi.ObjectType("java/lang/String"),
		jS,
	)
	if err != nil {
		return "", err
	}
	return k.javaInstance.JStringToString(jTokensString)
}

var jarPath = "tokenizers/dist/komoran"

func jarPaths() ([]string, error) {
	files, err := ioutil.ReadDir(jarPath)
	if err != nil {
		return nil, err
	}

	a := make([]string, len(files))
	for i, f := range files {
		a[i] = filepath.Join(jarPath, f.Name())
	}
	return a, nil
}
