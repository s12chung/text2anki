// Package main is the start point for text2anki
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/api"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

const host = "http://localhost"
const port = "3000"

var plog = logg.Default()

func configFromEnv() config.Config {
	c := config.Config{}

	c.Log = plog
	c.TxPool = api.TxPool{}
	if os.Getenv("TOKENIZER") == "komoran" {
		c.TokenizerType = config.TokenizerKomoran
	}
	if os.Getenv("DICTIONARY") == "koreanbasic" {
		c.DictionaryType = config.DictionaryKoreanBasic
	}
	c.StorageConfig = config.StorageConfig{
		StorageType: config.StorageLocalStore,
		LocalStoreConfig: config.LocalStoreConfig{
			Origin:        host + ":" + port,
			KeyBasePath:   "db/tmp/filestore",
			EncryptorPath: "db/tmp",
		},
	}
	return c
}

func main() {
	cli := false
	if cli {
		mainAgain()
		return
	}

	if err := run(); err != nil {
		plog.Error("main", logg.Err(err))
		os.Exit(-1)
	}
}

func run() error {
	if err := db.SetDB("db/data.sqlite3"); err != nil {
		return err
	}
	ctx := context.Background()
	routes := api.NewRoutes(ctx, configFromEnv())

	if err := routes.Setup(ctx); err != nil {
		return err
	}
	defer func() {
		if err := routes.Cleanup(); err != nil {
			plog.Error("main routes.Cleanup()", logg.Err(err))
		}
	}()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
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
	plog.Info("Server running on " + host + server.Addr)
	return server.Serve(ln)
}

func mainAgain() {
	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("usage: %v textStringFilename exportDir\n", os.Args[0]) //nolint:forbidigo // usage
		os.Exit(-1)
	}

	textStringFilename, exportDir := args[0], args[1]

	if err := runAgain(textStringFilename, exportDir); err != nil {
		plog.Error("main", logg.Err(err))
		os.Exit(-1)
	}
}

func runAgain(_, exportDir string) error {
	if err := anki.SetupDefaultConfig(); err != nil {
		return err
	}
	return exportFiles([]anki.Note{}, exportDir)
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
	synth := config.Synthesizer()
	for i := range notes {
		note := &notes[i]
		log := plog.With(slog.String("text", note.Text))

		speech, err := synth.TextToSpeech(context.Background(), note.Usage)
		if err != nil {
			log.Error("error creating audio for note", logg.Err(err))
		}
		if err = note.SetSound(speech, synth.SourceName()); err != nil {
			log.Error("error creating audio for note", logg.Err(err))
		}
	}
	return nil
}
