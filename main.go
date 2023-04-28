// Package main is the start point for text2anki
package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/cmd/prompt"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/synthesizers/azure"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizers/komoran"
	"github.com/s12chung/text2anki/pkg/util/stringclean"
)

var cleanSpeaker bool

func init() {
	flag.BoolVar(&cleanSpeaker, "clean-speaker", false, "clean 'speaker name:' from text")
	flag.Parse()
}

var tokenizer = func() tokenizers.Tokenizer {
	switch os.Getenv("TOKENIZER") {
	case "komoran":
		return komoran.New()
	default:
		return khaiii.New()
	}
}()
var parser = text.NewParser(text.Korean, text.English)
var dictionary = koreanbasic.New(koreanbasic.GetAPIKeyFromEnv())
var synth = azure.New(azure.GetAPIKeyFromEnv(), azure.EastUSRegion)

func main() {
	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("usage: %v textStringFilename exportDir\n", os.Args[0])
		os.Exit(-1)
	}

	textStringFilename, exportDir := args[0], args[1]

	if err := run(textStringFilename, exportDir); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run(textStringFilename, exportDir string) error {
	if err := anki.SetupDefaultConfig(); err != nil {
		return err
	}
	tokenizedTexts, err := tokenizeFile(textStringFilename)
	if err != nil {
		return err
	}
	notes, err := runUI(tokenizedTexts)
	if err != nil {
		return err
	}
	return exportFiles(notes, exportDir)
}

func tokenizeFile(filename string) ([]text.TokenizedText, error) {
	//nolint:gosec // required for binary to work
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	texts, err := parser.TextsFromString(string(fileBytes))
	if err != nil {
		bytes, _ := yaml.Marshal(texts)
		fmt.Println(string(bytes))
		return nil, err
	}
	texts = cleanTexts(texts)

	tokenizedTexts, err := text.TokenizeTexts(tokenizer, texts)
	if err != nil {
		return nil, err
	}
	return tokenizedTexts, err
}

func cleanTexts(texts []text.Text) []text.Text {
	if !cleanSpeaker {
		return texts
	}

	cleanedTexts := make([]text.Text, len(texts))
	for i, t := range texts {
		cleanedTexts[i] = text.Text{
			Text:        stringclean.Speaker(t.Text),
			Translation: stringclean.Speaker(t.Translation),
		}
	}
	return cleanedTexts
}

func runUI(tokenizedTexts []text.TokenizedText) ([]anki.Note, error) {
	notes, err := prompt.CreateCards(tokenizedTexts, dictionary)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func exportFiles(notes []anki.Note, exportDir string) error {
	if err := createAudio(notes); err != nil {
		return err
	}
	if err := os.Mkdir(exportDir, 0750); err != nil {
		return err
	}
	if err := anki.ExportFiles(notes, exportDir); err != nil {
		return err
	}
	return nil
}

func createAudio(notes []anki.Note) error {
	for i := range notes {
		note := &notes[i]
		speech, err := synth.TextToSpeech(note.Usage)
		if err != nil {
			slog.Error("error creating audio for note",
				slog.String("text", note.Text), slog.String("err", err.Error()))
		}
		if err = note.SetSound(speech, synth.SourceName()); err != nil {
			slog.Error("error creating audio for note",
				slog.String("text", note.Text), slog.String("err", err.Error()))
		}
	}
	return nil
}
