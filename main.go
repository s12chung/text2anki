// Package main is the start point for text2anki
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/api"
	"github.com/s12chung/text2anki/pkg/cmd/prompt"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

var cleanSpeaker bool
var cli bool

func init() {
	flag.BoolVar(&cli, "cli", false, "use cli")
	flag.BoolVar(&cleanSpeaker, "clean-speaker", false, "clean 'speaker name:' from text")
	flag.Parse()
}

func main() {
	if cli {
		mainAgain()
		return
	}

	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run() error {
	if err := api.DefaultRoutes.Setup(); err != nil {
		return err
	}
	defer func() {
		if err := api.DefaultRoutes.Cleanup(); err != nil {
			fmt.Println(err)
		}
	}()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Mount("/", api.DefaultRoutes.Router())

	server := http.Server{
		Addr:              ":3000",
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	slog.Info("Server running on http://localhost" + server.Addr)
	return server.Serve(ln)
}

func mainAgain() {
	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("usage: %v textStringFilename exportDir\n", os.Args[0])
		os.Exit(-1)
	}

	textStringFilename, exportDir := args[0], args[1]

	if err := runAgain(textStringFilename, exportDir); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func runAgain(textStringFilename, exportDir string) error {
	if err := api.DefaultRoutes.Setup(); err != nil {
		return err
	}
	defer func() {
		if err := api.DefaultRoutes.Cleanup(); err != nil {
			fmt.Println(err)
		}
	}()

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

func tokenizeFile(filename string) ([]db.TokenizedText, error) {
	//nolint:gosec // required for binary to work
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	texts, err := api.DefaultRoutes.TextTokenizer.Parser.TextsFromString(string(fileBytes))
	if err != nil {
		bytes, _ := yaml.Marshal(texts)
		fmt.Println(string(bytes))
		return nil, err
	}
	if cleanSpeaker {
		texts = text.CleanSpeaker(texts)
	}

	tokenizedTexts, err := api.DefaultRoutes.TextTokenizer.TokenizeTexts(texts)
	if err != nil {
		return nil, err
	}
	return tokenizedTexts, err
}

func runUI(tokenizedTexts []db.TokenizedText) ([]anki.Note, error) {
	notes, err := prompt.CreateCards(tokenizedTexts, api.DefaultRoutes.Dictionary)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func exportFiles(notes []anki.Note, exportDir string) error {
	if err := createAudio(notes); err != nil {
		return err
	}
	if err := os.Mkdir(exportDir, ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	if err := anki.ExportFiles(notes, exportDir); err != nil {
		return err
	}
	return nil
}

func createAudio(notes []anki.Note) error {
	synth := api.DefaultRoutes.Synthesizer
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
