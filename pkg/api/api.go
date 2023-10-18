// Package api contains the routes for the api
package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx/reqtxchi"
)

// Routes contains the routes used for the api
type Routes struct {
	CacheDir string
	Log      *slog.Logger

	UUIDGenerator storage.UUIDGenerator
	TxIntegrator  reqtx.Integrator[db.TxQs, config.TxMode]
	Storage       config.Storage

	Dictionary    dictionary.Dictionary
	TextTokenizer db.TextTokenizer

	ExtractorMap extractor.Map

	SoundSetter anki.SoundSetter
}

// NewRoutes is the routes used by the API
func NewRoutes(ctx context.Context, c config.Config) Routes {
	routes := Routes{
		CacheDir: config.CacheDir(c.CacheDir),
		Log:      c.Log,

		UUIDGenerator: config.UUIDGenerator(c.UUIDGenerator),
		TxIntegrator:  config.TxIntegrator(c.TxPool),
		Storage:       config.StorageFromConfig(c.StorageConfig, c.Log),

		Dictionary: config.Dictionary(c.DictionaryType),
		TextTokenizer: db.TextTokenizer{
			Parser:       config.Parser(),
			Tokenizer:    config.Tokenizer(ctx, c.TokenizerType, c.Log),
			Translator:   config.Translator(c.Translator),
			CleanSpeaker: true,
		},
		ExtractorMap: config.ExtractorMap(c.ExtractorMap),

		SoundSetter: config.SoundSetter(config.Synthesizer(c.Synthesizer)),
	}
	db.SetLog(c.Log)
	jhttp.SetLog(c.Log)
	db.SetDBStorage(routes.Storage.DBStorage)
	anki.SetConfig(config.Anki(c.AnkiCacheDir, c.Log))
	return routes
}

// Setup sets up the routes
func (rs Routes) Setup(ctx context.Context) error { return rs.TextTokenizer.Setup(ctx) }

// Cleanup cleans up the routes
func (rs Routes) Cleanup() error { return rs.TextTokenizer.Cleanup() }

// Router returns the router with all the routes set
func (rs Routes) Router() chi.Router {
	r := reqtxchi.NewRouter[db.TxQs, config.TxMode](chi.NewRouter(), rs.TxIntegrator, httpWrapper{})
	r.WithChi(func(r chi.Router) {
		r.NotFound(jhttp.ResponseWrap(rs.NotFound))
		r.MethodNotAllowed(jhttp.ResponseWrap(rs.NotAllowed))

		r.Route(config.StorageURLPath, func(r chi.Router) {
			r.Method(http.MethodGet, "/*", http.StripPrefix(config.StorageURLPath, rs.StorageGet()))
			r.Put("/*", responseJSONWrap(rs.StoragePut))
		})
	})

	r.Route("/sources", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
		r.Get("/", rs.SourcesIndex)
		r.Mode(txWritable).Post("/", rs.SourceCreate)

		r.Route("/{id}", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
			r.Get("/", rs.SourceGet)
			r.Mode(txWritable).Patch("/", rs.SourceUpdate)
			r.Mode(txWritable).Delete("/", rs.SourceDestroy)

			r.Route("/parts", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
				r.Mode(txWritable).Post("/", rs.PartCreate)
				r.Mode(txWritable).Post("/multi", rs.PartCreateMulti)
				r.Mode(txWritable).Patch("/{partIndex}", rs.PartUpdate)
				r.Mode(txWritable).Delete("/{partIndex}", rs.PartDestroy)
			})
		})

		r.Route("/pre_part_lists", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
			r.Post("/", rs.PrePartListCreate)
			r.Post("/sign", rs.PrePartListSign)
			r.Post("/verify", rs.PrePartListVerify)
			r.Route("/{id}", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
				r.Get("/", rs.PrePartListGet)
			})
		})
	})
	r.Route("/terms", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
		r.Get("/search", rs.TermsSearch)
	})
	r.Route("/notes", func(r reqtxchi.Router[db.TxQs, config.TxMode]) {
		r.Get("/", rs.NotesIndex)
		r.Mode(txWritable).Post("/", rs.NoteCreate)
		r.Mode(txWritable).Chi().Get("/download", rs.NotesDownload)
	})
	return r.Router
}

func responseJSONWrap(f jhttp.ResponseHandler) http.HandlerFunc {
	return jhttp.ResponseWrap(responseWrap(f))
}

type httpWrapper struct{}

func (h httpWrapper) RequestWrap(f jhttp.RequestHandler) jhttp.RequestHandler { return f }
func (h httpWrapper) ResponseWrap(f jhttp.ResponseHandler) jhttp.ResponseHandler {
	return responseWrap(f)
}

func responseWrap(f jhttp.ResponseHandler) jhttp.ResponseHandler { return prepareModelWrap(f) }
func prepareModelWrap(f jhttp.ResponseHandler) jhttp.ResponseHandler {
	return func(r *http.Request) (any, *jhttp.HTTPError) {
		model, httpErr := f(r)
		if httpErr != nil {
			return model, httpErr
		}
		if err := httptyped.PrepareModel(model); err != nil {
			return model, jhttp.Error(http.StatusInternalServerError, err)
		}
		return model, nil
	}
}

func (rs Routes) runOr500(r *http.Request, f func(r *http.Request, tx db.TxQs) error) *jhttp.HTTPError {
	_, httpErr := rs.TxIntegrator.ResponseWrap(func(r *http.Request, tx db.TxQs) (any, *jhttp.HTTPError) {
		return jhttp.ReturnModelOr500(func() (any, error) { return nil, f(r, tx) })
	})(r)
	return httpErr
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
