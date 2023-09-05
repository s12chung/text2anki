// Package reqtxchi wraps chi.Router functions to adhere to jhttp.RequestHandler and the custom ResponseHandler
// nolint: revive
package reqtxchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
)

type ResponseHandler[T reqtx.Tx] func(r *http.Request, tx T) (any, *jhttp.HTTPError)

// HTTPWrapper provides wrapping functions that Mux uses
type HTTPWrapper interface {
	RequestWrap(f jhttp.RequestHandler) jhttp.RequestHandler
	ResponseWrap(f jhttp.ResponseHandler) jhttp.ResponseHandler
}

// Router represents a ResponseHandler router
//
//nolint:interfacebloat // just following chi.Router
type Router[T reqtx.Tx] interface {
	Use(middlewares ...jhttp.RequestHandler)
	With(middlewares ...jhttp.RequestHandler) Router[T]

	Group(fn func(r Router[T])) Router[T]
	Route(pattern string, fn func(r Router[T])) Router[T]

	HandleFunc(pattern string, f ResponseHandler[T])
	MethodFunc(method, pattern string, f ResponseHandler[T])

	Connect(pattern string, f ResponseHandler[T])
	Delete(pattern string, f ResponseHandler[T])
	Get(pattern string, f ResponseHandler[T])
	Head(pattern string, f ResponseHandler[T])
	Options(pattern string, f ResponseHandler[T])
	Patch(pattern string, f ResponseHandler[T])
	Post(pattern string, f ResponseHandler[T])
	Put(pattern string, f ResponseHandler[T])
	Trace(pattern string, f ResponseHandler[T])
}

// Mux returns a new Mux that wraps chi.Router functions
type Mux[T reqtx.Tx] struct {
	Router  chi.Router
	wrapper HTTPWrapper
}

// NewRouter returns a new Mux
func NewRouter[T reqtx.Tx](r chi.Router, wrapper HTTPWrapper) Mux[T] {
	return Mux[T]{Router: r, wrapper: wrapper}
}

func (m Mux[T]) router(r chi.Router) Mux[T] {
	dup := m
	dup.Router = r
	return dup
}

func (m Mux[T]) requestWrap(f jhttp.RequestHandler) func(http.Handler) http.Handler {
	return jhttp.RequestWrap(m.wrapper.RequestWrap(f))
}
func (m Mux[T]) responseWrap(f ResponseHandler[T]) http.HandlerFunc {
	return jhttp.ResponseWrap(m.wrapper.ResponseWrap(ResponseWrap(f)))
}

func ResponseWrap[T reqtx.Tx](f ResponseHandler[T]) jhttp.ResponseHandler {
	return func(r *http.Request) (any, *jhttp.HTTPError) {
		tx, err := reqtx.ContextTx[T](r)
		if err != nil {
			return nil, jhttp.Error(http.StatusInternalServerError, err)
		}
		return f(r, tx)
	}
}

func (m Mux[T]) Use(middlewares ...jhttp.RequestHandler) {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = m.requestWrap(middleware)
	}
	m.Router.Use(wrapped...)
}
func (m Mux[T]) With(middlewares ...jhttp.RequestHandler) Router[T] {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = m.requestWrap(middleware)
	}
	return m.router(m.Router.With(wrapped...))
}

func (m Mux[T]) Group(fn func(r Router[T])) Router[T] {
	r := m.Router.Group(func(r chi.Router) {
		fn(m.router(r))
	})
	return m.router(r)
}
func (m Mux[T]) Route(pattern string, fn func(r Router[T])) Router[T] {
	r := m.Router.Route(pattern, func(r chi.Router) {
		fn(m.router(r))
	})
	return m.router(r)
}

func (m Mux[T]) HandleFunc(pattern string, f ResponseHandler[T]) {
	m.Router.HandleFunc(pattern, m.responseWrap(f))
}
func (m Mux[T]) MethodFunc(method, pattern string, f ResponseHandler[T]) {
	m.Router.MethodFunc(method, pattern, m.responseWrap(f))
}

func (m Mux[T]) Connect(pattern string, f ResponseHandler[T]) {
	m.Router.Connect(pattern, m.responseWrap(f))
}
func (m Mux[T]) Delete(pattern string, f ResponseHandler[T]) {
	m.Router.Delete(pattern, m.responseWrap(f))
}
func (m Mux[T]) Get(pattern string, f ResponseHandler[T]) {
	m.Router.Get(pattern, m.responseWrap(f))
}
func (m Mux[T]) Head(pattern string, f ResponseHandler[T]) {
	m.Router.Head(pattern, m.responseWrap(f))
}
func (m Mux[T]) Options(pattern string, f ResponseHandler[T]) {
	m.Router.Options(pattern, m.responseWrap(f))
}
func (m Mux[T]) Patch(pattern string, f ResponseHandler[T]) {
	m.Router.Patch(pattern, m.responseWrap(f))
}
func (m Mux[T]) Post(pattern string, f ResponseHandler[T]) {
	m.Router.Post(pattern, m.responseWrap(f))
}
func (m Mux[T]) Put(pattern string, f ResponseHandler[T]) {
	m.Router.Put(pattern, m.responseWrap(f))
}
func (m Mux[T]) Trace(pattern string, f ResponseHandler[T]) {
	m.Router.Trace(pattern, m.responseWrap(f))
}
