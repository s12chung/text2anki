// Package api contains the routes for the api
package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/synthesizer"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/jhttp/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp/jchi"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
)

// Routes contains the routes used for the api
type Routes struct {
	TxIntegrator reqtx.Integrator

	Dictionary    dictionary.Dictionary
	Synthesizer   synthesizer.Synthesizer
	TextTokenizer db.TextTokenizer

	Storage      config.Storage
	ExtractorMap extractor.Map
}

// NewRoutes is the routes used by the API
func NewRoutes(c config.Config) Routes {
	routes := Routes{
		TxIntegrator: config.TxIntegrator(c.TxPool),

		Dictionary:  config.Dictionary(c.DictionaryType),
		Synthesizer: config.Synthesizer(),
		TextTokenizer: db.TextTokenizer{
			Parser:       config.Parser(),
			Tokenizer:    config.Tokenizer(c.TokenizerType),
			CleanSpeaker: true,
		},

		Storage:      config.StorageFromConfig(c.StorageConfig),
		ExtractorMap: config.ExtractorMap(c.ExtractorMap),
	}
	db.SetDBStorage(routes.Storage.DBStorage)
	return routes
}

// Setup sets up the routes
func (rs Routes) Setup() error { return rs.TextTokenizer.Setup() }

// Cleanup cleans up the routes
func (rs Routes) Cleanup() error { return rs.TextTokenizer.Cleanup() }

// Router returns the router with all the routes set
func (rs Routes) Router() chi.Router {
	var r chi.Router = chi.NewRouter()
	r.NotFound(jhttp.ResponseJSONWrap(rs.NotFound))
	r.MethodNotAllowed(jhttp.ResponseJSONWrap(rs.NotAllowed))

	r.Route(config.StorageURLPath, func(r chi.Router) {
		r.Method(http.MethodGet, "/*", http.StripPrefix(config.StorageURLPath, rs.StorageGet()))
		r.Put("/*", responseJSONWrap(rs.StoragePut))
	})

	r.Mount("/", rs.txRouter())
	return r
}

func (rs Routes) txRouter() chi.Router {
	r := jchi.NewRouter(chi.NewRouter(), httpWrapper{})
	r.Router.Use(jhttp.RequestWrap(rs.TxIntegrator.SetTxContext))

	r.Route("/sources", func(r jchi.Router) {
		r.Get("/", rs.SourceIndex)
		r.Post("/", rs.SourceCreate)

		r.Route("/{sourceID}", func(r jchi.Router) {
			r.Use(rs.SourceCtx)
			r.Get("/", rs.SourceGet)
			r.Patch("/", rs.SourceUpdate)
			r.Delete("/", rs.SourceDestroy)
		})

		r.Route("/pre_part_lists", func(r jchi.Router) {
			r.Post("/", rs.PrePartListCreate)
			r.Post("/sign", rs.PrePartListSign)
			r.Post("/verify", rs.PrePartListVerify)
			r.Route("/{prePartListID}", func(r jchi.Router) {
				r.Get("/", rs.PrePartListGet)
			})
		})
	})
	r.Route("/terms", func(r jchi.Router) {
		r.Get("/search", rs.TermsSearch)
	})
	r.Route("/notes", func(r jchi.Router) {
		r.Post("/", rs.NoteCreate)
	})
	return r.Router
}

func responseJSONWrap(f jhttp.ResponseJSONWrapFunc) http.HandlerFunc {
	return jhttp.ResponseJSONWrap(responseWrap(f))
}

func responseWrap(f jhttp.ResponseJSONWrapFunc) jhttp.ResponseJSONWrapFunc {
	return httptyped.TypedWrap(f)
}

type httpWrapper struct{}

func (h httpWrapper) RequestWrap(f jhttp.RequestWrapFunc) jhttp.RequestWrapFunc {
	return reqtx.TxRollbackRequestWrap(f)
}
func (h httpWrapper) ResponseWrap(f jhttp.ResponseJSONWrapFunc) jhttp.ResponseJSONWrapFunc {
	return reqtx.TxFinalizeWrap(responseWrap(f))
}

// NotFound is the route handler for not matching pattern routes
func (rs Routes) NotFound(r *http.Request) (any, *jhttp.HTTPError) {
	return nil, jhttp.Error(http.StatusNotFound, fmt.Errorf("request URL, %v, does not match any route", r.URL.String()))
}

// NotAllowed is the router handler for method not handled for the pattern
func (rs Routes) NotAllowed(r *http.Request) (any, *jhttp.HTTPError) {
	return nil, jhttp.Error(http.StatusMethodNotAllowed,
		fmt.Errorf("the method, %v (at %v), is not allowed with at this URL", r.Method, r.URL.String()))
}
