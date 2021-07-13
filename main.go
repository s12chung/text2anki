package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/tokenizer/komoran"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run() error {
	tokenizer := komoran.NewKomoran()
	if err := tokenizer.Setup(); err != nil {
		return err
	}
	tokens, err := tokenizer.GetTokens("대한민국은 민주공화국이다.")
	if err != nil {
		return err
	}

	tokensJSON, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(tokensJSON))

	dict := koreanbasic.NewKoreanBasic(os.Getenv("KOREAN_BASIC_API_KEY"))
	terms, err := dict.Search(tokens[0].Morph)
	if err != nil {
		return err
	}

	termsJSON, err := json.MarshalIndent(terms, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(termsJSON))

	return nil
}
