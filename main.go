package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/s12chung/text2anki/cmd/prompt"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/text"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run() error {
	// TODO: clean this up
	unixTimeMs := time.Now().UnixNano() / int64(time.Millisecond)
	dir := fmt.Sprintf("tmp/%v", unixTimeMs)
	err := os.Mkdir(dir, 0750)
	if err != nil {
		return err
	}
	textStringFilename, exportDir := "tmp/in.txt", dir

	tokenizedTexts, err := tokenizeTexts(textStringFilename)
	if err != nil {
		return err
	}

	notes, err := runUI(tokenizedTexts)
	if err != nil {
		return err
	}

	if err := anki.ExportFiles(notes, exportDir); err != nil {
		return err
	}
	return nil
}

func tokenizeTexts(textStringFilename string) ([]app.TokenizedText, error) {
	textString, err := readTextString(textStringFilename)
	if err != nil {
		return nil, err
	}
	texts := text.TextsFromString(textString)
	tokenizedTexts, err := app.TokenizeTexts(texts)
	if err != nil {
		return nil, err
	}
	return tokenizedTexts, err
}

func readTextString(filename string) (string, error) {
	//nolint:gosec // required for binary to work
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func runUI(tokenizedTexts []app.TokenizedText) ([]anki.Note, error) {
	dictionary := koreanbasic.NewKoreanBasic(koreanbasic.GetAPIKeyFromEnv())
	notes, err := prompt.Revolve(tokenizedTexts, dictionary)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
