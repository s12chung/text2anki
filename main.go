// Package main is the start point for text2anki
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/api"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

var routes = api.NewRoutes(configFromEnv())

const host = "http://localhost"
const port = "3000"

func configFromEnv() api.Config {
	config := api.Config{}
	if os.Getenv("TOKENIZER") == "komoran" {
		config.TokenizerType = api.TokenizerKomoran
	}
	if os.Getenv("DICTIONARY") == "koreanbasic" {
		config.DictionaryType = api.DictionaryKoreanBasic
	}
	config.SignerConfig = api.SignerConfig{
		SignerType: api.SignerLocalStore,
		LocalStoreConfig: api.LocalStoreConfig{
			Origin:   host + ":" + port,
			BaseBath: "db/tmp/filestore",
			KeyPath:  "db/tmp",
		},
	}
	return config
}

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
	if err := routes.Setup(); err != nil {
		return err
	}
	defer func() {
		if err := routes.Cleanup(); err != nil {
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
	r.Mount("/", routes.Router())

	server := http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	slog.Info("Server running on " + host + server.Addr)
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
	if err := routes.Setup(); err != nil {
		return err
	}
	defer func() {
		if err := routes.Cleanup(); err != nil {
			fmt.Println(err)
		}
	}()

	if err := anki.SetupDefaultConfig(); err != nil {
		return err
	}
	_, err := tokenizeFile(textStringFilename)
	if err != nil {
		return err
	}
	return exportFiles([]anki.Note{}, exportDir)
}

func tokenizeFile(filename string) ([]db.TokenizedText, error) {
	//nolint:gosec // required for binary to work
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	split := strings.Split(string(fileBytes), "===")
	if len(split) == 1 {
		split = append(split, "")
	}
	texts, err := routes.TextTokenizer.Parser.Texts(split[0], split[1])
	if err != nil {
		bytes, _ := yaml.Marshal(texts)
		fmt.Println(string(bytes))
		return nil, err
	}
	if cleanSpeaker {
		texts = text.CleanSpeaker(texts)
	}

	tokenizedTexts, err := routes.TextTokenizer.TokenizeTexts(texts)
	if err != nil {
		return nil, err
	}
	return tokenizedTexts, err
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
	synth := routes.Synthesizer
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
