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
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/synthesizer"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

// Routes contains the routes used for the api
type Routes struct {
	Dictionary    dictionary.Dictionary
	Synthesizer   synthesizer.Synthesizer
	TextTokenizer db.TextTokenizer
	Storage       config.Storage
	ExtractorMap  extractor.Map
}

// NewRoutes is the routes used by the API
func NewRoutes(c config.Config) Routes {
	routes := Routes{
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
		r.Get("/", responseWrap(rs.SourceIndex))
		r.Post("/", responseWrap(rs.SourceCreate))

		r.Route("/{sourceID}", func(r chi.Router) {
			r.Use(httputil.RequestWrap(SourceCtx))
			r.Get("/", responseWrap(rs.SourceGet))
			r.Patch("/", responseWrap(rs.SourceUpdate))
			r.Delete("/", responseWrap(rs.SourceDestroy))
		})

		r.Route("/pre_part_lists", func(r chi.Router) {
			r.Post("/", responseWrap(rs.PrePartListCreate))
			r.Post("/sign", responseWrap(rs.PrePartListSign))
			r.Post("/verify", responseWrap(rs.PrePartListVerify))
			r.Route("/{prePartListID}", func(r chi.Router) {
				r.Get("/", responseWrap(rs.PrePartListGet))
			})
		})
	})
	r.Route("/terms", func(r chi.Router) {
		r.Get("/search", responseWrap(rs.TermsSearch))
	})
	r.Route("/notes", func(r chi.Router) {
		r.Post("/", responseWrap(rs.NoteCreate))
	})
	r.Route(config.StorageURLPath, func(r chi.Router) {
		r.Method(http.MethodGet, "/*", http.StripPrefix(config.StorageURLPath, rs.StorageGet()))
		r.Put("/*", responseWrap(rs.StoragePut))
	})
	r.NotFound(httputil.RespondJSONWrap(rs.NotFound))
	r.MethodNotAllowed(httputil.RespondJSONWrap(rs.NotAllowed))
	return r
}

func responseWrap(f httputil.RespondJSONWrapFunc) http.HandlerFunc {
	return httputil.RespondJSONWrap(httptyped.TypedWrap(f))
}

// NotFound is the route handler for not matching pattern routes
func (rs Routes) NotFound(r *http.Request) (any, *httputil.HTTPError) {
	return nil, httputil.Error(http.StatusNotFound, fmt.Errorf("request URL, %v, does not match any route", r.URL.String()))
}

// NotAllowed is the router handler for method not handled for the pattern
func (rs Routes) NotAllowed(r *http.Request) (any, *httputil.HTTPError) {
	return nil, httputil.Error(http.StatusMethodNotAllowed,
		fmt.Errorf("the method, %v (at %v), is not allowed with at this URL", r.Method, r.URL.String()))
}

func extractAndValidate(r *http.Request, req any) *httputil.HTTPError {
	if httpError := httputil.ExtractJSON(r, req); httpError != nil {
		return httpError
	}
	result := firm.Validate(req)
	if !result.IsValid() {
		return httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf(result.ErrorMap().String()))
	}
	return nil
}
