// Package api contains the routes for the api
package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/synthesizers"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

// Routes contains the routes used for the api
type Routes struct {
	Dictionary    dictionary.Dictionary
	Synthesizer   synthesizers.Synthesizer
	TextTokenizer db.TextTokenizer
	Storage       Storage
}

// Storage contains the Route's storage setup
type Storage struct {
	Signer storage.Signer
	Storer storage.Storer
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
		r.Get("/", httptyped.RespondTypedJSONWrap(rs.SourceIndex))
		r.Post("/", httptyped.RespondTypedJSONWrap(rs.SourceCreate))

		r.Route("/{sourceID}", func(r chi.Router) {
			r.Use(httputil.RequestWrap(SourceCtx))
			r.Get("/", httptyped.RespondTypedJSONWrap(rs.SourceGet))
			r.Patch("/", httptyped.RespondTypedJSONWrap(rs.SourceUpdate))
			r.Delete("/", httptyped.RespondTypedJSONWrap(rs.SourceDestroy))
		})

		r.Route("/pre_part_lists", func(r chi.Router) {
			r.Post("/sign", httptyped.RespondTypedJSONWrap(rs.PrePartListSign))
			r.Route("/{prePartListID}", func(r chi.Router) {
				r.Get("/", httptyped.RespondTypedJSONWrap(rs.PrePartListGet))
			})
		})
	})
	r.Route("/terms", func(r chi.Router) {
		r.Get("/search", httptyped.RespondTypedJSONWrap(rs.TermsSearch))
	})
	r.Route("/notes", func(r chi.Router) {
		r.Post("/", httptyped.RespondTypedJSONWrap(rs.NoteCreate))
	})
	r.Route(storageURLPath, func(r chi.Router) {
		r.Method(http.MethodGet, "/*", http.StripPrefix(storageURLPath, rs.StorageGet()))
		r.Put("/*", httptyped.RespondTypedJSONWrap(rs.StoragePut))
	})
	r.NotFound(httputil.RespondJSONWrap(rs.NotFound))
	r.MethodNotAllowed(httputil.RespondJSONWrap(rs.NotAllowed))
	return r
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

type signFieldConfig struct {
	Name      string
	ValidExts map[string]bool
}

func signFieldIfExists(builder storage.SignPutBuilder, config signFieldConfig, ext string) (*storage.PreSignedHTTPRequest, *httputil.HTTPError) {
	if ext == "" {
		return nil, nil
	}
	if !config.ValidExts[ext] {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("%v is not a valid file extension", ext))
	}
	signedReq, err := builder.Field(config.Name).Sign(ext)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return &signedReq, nil
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
