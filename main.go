package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/s12chung/text2anki/cmd/prompt"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/text"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v textStringFilename exportDir\n", os.Args[0])
		os.Exit(-1)
	}

	textStringFilename, exportDir := os.Args[1], os.Args[2]

	if err := run(textStringFilename, exportDir); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run(textStringFilename, exportDir string) error {
	tokenizedTexts, err := tokenizeTexts(textStringFilename)
	if err != nil {
		return err
	}

	notes, err := runUI(tokenizedTexts)
	if err != nil {
		return err
	}

	return exportFiles(notes, exportDir)
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
	notes, err := prompt.CreateCards(tokenizedTexts, dictionary)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func exportFiles(notes []anki.Note, exportDir string) error {
	err := os.Mkdir(exportDir, 0750)
	if err != nil {
		return err
	}

	if err := anki.ExportFiles(notes, exportDir); err != nil {
		return err
	}
	return nil
}
