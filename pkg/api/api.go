// Package api contains the routes for the api
package api

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/dictionary/krdict"
	"github.com/s12chung/text2anki/pkg/synthesizers"
	"github.com/s12chung/text2anki/pkg/synthesizers/azure"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/tokenizers/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizers/komoran"
	"github.com/s12chung/text2anki/pkg/util/httputil"
)

// DefaultRoutes is the routes used by the API
var DefaultRoutes = Routes{
	Dictionary:  DefaultDictionary(),
	Synthesizer: DefaultSynthesizer(),
	TextTokenizer: db.TextTokenizer{
		Parser:       DefaultParser(),
		Tokenizer:    DefaultTokenizer(),
		CleanSpeaker: true,
	},
}

// DefaultParser returns the default Parser
func DefaultParser() text.Parser {
	return text.NewParser(text.Korean, text.English)
}

// DefaultSynthesizer returns the default Synthesizer
func DefaultSynthesizer() synthesizers.Synthesizer {
	return azure.New(azure.GetAPIKeyFromEnv(), azure.EastUSRegion)
}

// DefaultTokenizer returns the default Tokenizer
func DefaultTokenizer() tokenizers.Tokenizer {
	switch os.Getenv("TOKENIZER") {
	case "komoran":
		return komoran.New()
	default:
		return khaiii.New()
	}
}

// DefaultDictionary returns the default Dictionary
func DefaultDictionary() dictionary.Dictionary {
	switch os.Getenv("DICTIONARY") {
	case "koreanbasic":
		return koreanbasic.New(koreanbasic.GetAPIKeyFromEnv())
	default:
		if err := db.SetDB("db/data.sqlite3"); err != nil {
			fmt.Println("failure to SetDB()\n", err)
			os.Exit(-1)
		}
		return krdict.New(db.DB())
	}
}

// Routes contains the routes used for the api
type Routes struct {
	Dictionary    dictionary.Dictionary
	Synthesizer   synthesizers.Synthesizer
	TextTokenizer db.TextTokenizer
}

// Setup sets up the routes
func (rs Routes) Setup() error {
	return rs.TextTokenizer.Setup()
}

// Cleanup cleans up the routes
func (rs Routes) Cleanup() error {
	return rs.TextTokenizer.Cleanup()
}

// Router returns the router with all the routes set
func (rs Routes) Router() chi.Router {
	r := chi.NewRouter()
	r.Route("/sources", func(r chi.Router) {
		r.Get("/", httputil.RespondJSONWrap(rs.SourceList))
		r.Post("/", httputil.RespondJSONWrap(rs.SourceCreate))

		r.Route("/{sourceID}", func(r chi.Router) {
			r.Use(httputil.RequestWrap(SourceCtx))
			r.Get("/", httputil.RespondJSONWrap(rs.SourceGet))
			r.Delete("/", httputil.RespondJSONWrap(rs.SourceDestroy))
		})
	})
	return r
}
