package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/app"
	"github.com/s12chung/text2anki/pkg/cmd/prompt"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/stringclean"
	"github.com/s12chung/text2anki/pkg/synthesizers/azure"
	"github.com/s12chung/text2anki/pkg/text"
)

var cleanSpeaker bool

func init() {
	flag.BoolVar(&cleanSpeaker, "clean-speaker", false, "clean 'speaker name:' from text")
	flag.Parse()
}

func main() {
	args := flag.Args()
	fmt.Println(args)

	if len(args) != 2 {
		fmt.Printf("Usage: %v textStringFilename exportDir\n", os.Args[0])
		os.Exit(-1)
	}

	textStringFilename, exportDir := args[0], args[1]

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
		var bytes []byte
		bytes, err = yaml.Marshal(texts)
		fmt.Println(string(bytes))
		return nil, err
	}

	cleanedTexts := make([]text.Text, len(texts))
	for i, t := range texts {
		cleanedTexts[i] = text.Text{
			Text:        cleanText(t.Text),
			Translation: cleanText(t.Translation),
		}
	}

	tokenizedTexts, err := app.TokenizeTexts(cleanedTexts)
	if err != nil {
		return nil, err
	}
	return tokenizedTexts, err
}

func cleanText(s string) string {
	if cleanSpeaker {
		s = stringclean.Speaker(s)
	}
	return s
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
