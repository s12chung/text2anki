// Package main is the start point for text2anki
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/api"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

const host = "http://localhost"
const port = "3000"

var plog = logg.Default() //nolint:forbidigo // main package

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
	if err := run(); err != nil {
		plog.Error("main", logg.Err(err))
		os.Exit(-1)
	}
}

func run() error {
	if err := db.SetDB("db/data.sqlite3"); err != nil {
		return err
	}
	ctx := context.Background() //nolint:forbidigo // this is main
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
