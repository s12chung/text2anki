// Package komoran is an interface to the Komoran Korean tokenizer
package komoran

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/pkg/java"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"tekao.net/jnigi"
)

// New returns a Komoran Korean tokenizer
func New() tokenizers.Tokenizer {
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
func (k *Komoran) GetTokens(s string) ([]tokenizers.Token, error) {
	if !k.javaInstance.IsSetup() {
		return nil, &tokenizers.NotSetupError{}
	}
	tokensJSON, err := k.callGetTokens(s)
	if err != nil {
		return nil, err
	}
	resp := &response{}
	if err := json.Unmarshal([]byte(tokensJSON), resp); err != nil {
		return nil, err
	}
	return resp.toTokenizerTokens(), nil
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
	k.javaInstance.Env.DeleteLocalRef(jS)

	ret, err := k.javaInstance.JStringToString(jTokensString)
	if err != nil {
		return "", err
	}
	k.javaInstance.Env.DeleteLocalRef(jTokensString.(*jnigi.ObjectRef))
	return ret, nil
}

var jarPath = "tokenizers/dist/komoran"

func jarPaths() ([]string, error) {
	files, err := os.ReadDir(jarPath)
	if err != nil {
		return nil, err
	}

	a := make([]string, len(files))
	for i, f := range files {
		a[i] = filepath.Join(jarPath, f.Name())
	}
	return a, nil
}
