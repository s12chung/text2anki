package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/cmd/prompt"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/synthesizers/azure"
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
	err := anki.SetupDefaultConfig()
	if err != nil {
		return err
	}

	tokenizedTexts, err := tokenizeTexts(textStringFilename)
	if err != nil {
		return err
	}

	notes, err := runUI(tokenizedTexts)
	if err != nil {
		return err
	}
	if err = createAudio(notes); err != nil {
		return err
	}

	return exportFiles(notes, exportDir)
}

func tokenizeTexts(textStringFilename string) ([]app.TokenizedText, error) {
	textString, err := readTextString(textStringFilename)
	if err != nil {
		return nil, err
	}

	parser := text.NewParser(text.Korean, text.English)
	texts, err := parser.TextsFromString(textString)
	if err != nil {
		bytes, err := yaml.Marshal(texts)
		fmt.Println(string(bytes))
		return nil, err
	}

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
	dictionary := koreanbasic.New(koreanbasic.GetAPIKeyFromEnv())
	notes, err := prompt.CreateCards(tokenizedTexts, dictionary)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func createAudio(notes []anki.Note) error {
	synth := azure.New(azure.GetAPIKeyFromEnv(), azure.EastUSRegion)
	for i := range notes {
		note := &notes[i]
		speech, err := synth.TextToSpeech(note.Usage)
		if err != nil {
			fmt.Printf("error creating audio for note (%v): %v\n", note.Text, err)
		}
		if err = note.SetSound(speech, synth.SourceName()); err != nil {
			fmt.Printf("error setting audio for note (%v): %v\n", note.Text, err)
		}
		time.Sleep(time.Second)
	}
	return nil
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
