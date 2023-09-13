// Package reqtxchi wraps chi.Router functions to adhere to jhttp.RequestHandler and the custom reqtx.ResponseHandler
// nolint: revive
package reqtxchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
)

// HTTPWrapper provides wrapping functions that Mux uses
type HTTPWrapper interface {
	RequestWrap(f jhttp.RequestHandler) jhttp.RequestHandler
	ResponseWrap(f jhttp.ResponseHandler) jhttp.ResponseHandler
}

// Router represents a reqtx.ResponseHandler router
//
//nolint:interfacebloat // just following chi.Router
type Router[T reqtx.Tx, Mode ~int] interface {
	Use(middlewares ...jhttp.RequestHandler)
	With(middlewares ...jhttp.RequestHandler) Router[T, Mode]

	Group(fn func(r Router[T, Mode])) Router[T, Mode]
	Route(pattern string, fn func(r Router[T, Mode])) Router[T, Mode]

	HandleFunc(pattern string, f reqtx.ResponseHandler[T])
	MethodFunc(method, pattern string, f reqtx.ResponseHandler[T])

	Connect(pattern string, f reqtx.ResponseHandler[T])
	Delete(pattern string, f reqtx.ResponseHandler[T])
	Get(pattern string, f reqtx.ResponseHandler[T])
	Head(pattern string, f reqtx.ResponseHandler[T])
	Options(pattern string, f reqtx.ResponseHandler[T])
	Patch(pattern string, f reqtx.ResponseHandler[T])
	Post(pattern string, f reqtx.ResponseHandler[T])
	Put(pattern string, f reqtx.ResponseHandler[T])
	Trace(pattern string, f reqtx.ResponseHandler[T])

	Mode(mode Mode) Router[T, Mode]
	Chi() chi.Router
	WithChi(fn func(r chi.Router))
}

// Mux returns a new Mux that wraps chi.Router functions
type Mux[T reqtx.Tx, Mode ~int] struct {
	Router     chi.Router
	integrator reqtx.Integrator[T, Mode]
	wrapper    HTTPWrapper
}

// NewRouter returns a new Mux
func NewRouter[T reqtx.Tx, Mode ~int](r chi.Router, integrator reqtx.Integrator[T, Mode], wrapper HTTPWrapper) Mux[T, Mode] {
	return Mux[T, Mode]{Router: r, integrator: integrator, wrapper: wrapper}
}

// Mode sets the transaction mode
func (m Mux[T, Mode]) Mode(mode Mode) Router[T, Mode] {
	return m.With(func(r *http.Request) (*http.Request, *jhttp.HTTPError) {
		return m.integrator.SetTxModeContext(r, mode), nil
	})
}

// Chi returns the chi.Router
func (m Mux[T, Mode]) Chi() chi.Router { return m.Router }

// WithChi creates a grouping with the chi.Router
func (m Mux[T, Mode]) WithChi(fn func(r chi.Router)) { fn(m.Router) }

func (m Mux[T, Mode]) requestWrap(f jhttp.RequestHandler) func(http.Handler) http.Handler {
	return jhttp.RequestWrap(m.wrapper.RequestWrap(f))
}
func (m Mux[T, Mode]) responseWrap(f reqtx.ResponseHandler[T]) http.HandlerFunc {
	return jhttp.ResponseWrap(m.wrapper.ResponseWrap(m.integrator.ResponseWrap(f)))
}

func (m Mux[T, Mode]) router(r chi.Router) Router[T, Mode] {
	dup := m
	dup.Router = r
	return dup
}

func (m Mux[T, Mode]) Use(middlewares ...jhttp.RequestHandler) {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = m.requestWrap(middleware)
	}
	m.Router.Use(wrapped...)
}
func (m Mux[T, Mode]) With(middlewares ...jhttp.RequestHandler) Router[T, Mode] {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = m.requestWrap(middleware)
	}
	return m.router(m.Router.With(wrapped...))
}

func (m Mux[T, Mode]) Group(fn func(r Router[T, Mode])) Router[T, Mode] {
	r := m.Router.Group(func(r chi.Router) {
		fn(m.router(r))
	})
	return m.router(r)
}
func (m Mux[T, Mode]) Route(pattern string, fn func(r Router[T, Mode])) Router[T, Mode] {
	r := m.Router.Route(pattern, func(r chi.Router) {
		fn(m.router(r))
	})
	return m.router(r)
}

func (m Mux[T, Mode]) HandleFunc(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.HandleFunc(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) MethodFunc(method, pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.MethodFunc(method, pattern, m.responseWrap(f))
}

func (m Mux[T, Mode]) Connect(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Connect(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Delete(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Delete(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Get(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Get(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Head(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Head(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Options(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Options(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Patch(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Patch(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Post(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Post(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Put(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Put(pattern, m.responseWrap(f))
}
func (m Mux[T, Mode]) Trace(pattern string, f reqtx.ResponseHandler[T]) {
	m.Router.Trace(pattern, m.responseWrap(f))
}
