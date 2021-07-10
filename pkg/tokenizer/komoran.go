package tokenizer

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"tekao.net/jnigi"
)

// NewKomoran returns a Komoran Korean tokenizer
func NewKomoran() Tokenizer {
	return &Komoran{}
}

// Komoran is a Korean Tokenizer in java
type Komoran struct {
	javaInstance
}

const jarPath = "tokenizers/build/komoran"

// Setup setups up the JVM for Komoran to run
func (k *Komoran) Setup() error {
	classPathArray, err := jarPaths()
	if err != nil {
		return err
	}
	return k.javaInstance.setup(strings.Join(append(classPathArray, ""), ":"))
}

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

type komoranTokens struct {
	TokenList []Token `json:"tokenList"`
}

// GetTokens returns the part of speech tokens of the given string
func (k *Komoran) GetTokens(s string) ([]Token, error) {
	if !k.IsSetup() {
		return nil, &NotSetupError{}
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

func (k *Komoran) callGetTokens(s string) (string, error) {
	jS, err := k.env.NewObject("java/lang/String", []byte(s))
	if err != nil {
		return "", err
	}

	jTokensString, err := k.env.CallStaticMethod("text2anki/tokenizer/komoran/Tokenizer",
		"getTokens",
		jnigi.ObjectType("java/lang/String"),
		jS,
	)
	if err != nil {
		return "", err
	}
	return k.jStringToString(jTokensString)
}
