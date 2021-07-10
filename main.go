package main

import (
	"fmt"
	"os"

	"github.com/s12chung/text2anki/pkg/tokenizer"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run() error {
	tokenizer := tokenizer.NewKomoran()
	if err := tokenizer.Setup(); err != nil {
		return err
	}
	tokens, err := tokenizer.GetTokens()
	if err != nil {
		return err
	}
	fmt.Println(tokens)
	return nil
}
