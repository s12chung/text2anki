// Package jchi wraps chi.Router functions to adhere to jhttp.ResponseHandler and jhttp.RequestHandler
// nolint: revive
package jchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

// HTTPWrapper provides wrapping functions that Mux uses
type HTTPWrapper interface {
	RequestWrap(f jhttp.RequestHandler) jhttp.RequestHandler
	ResponseWrap(f jhttp.ResponseHandler) jhttp.ResponseHandler
}

// Router represents a jchi router
//
//nolint:interfacebloat // just following chi.Router
type Router interface {
	Use(middlewares ...jhttp.RequestHandler)
	With(middlewares ...jhttp.RequestHandler) Router

	Group(fn func(r Router)) Router
	Route(pattern string, fn func(r Router)) Router

	HandleFunc(pattern string, f jhttp.ResponseHandler)
	MethodFunc(method, pattern string, f jhttp.ResponseHandler)

	Connect(pattern string, f jhttp.ResponseHandler)
	Delete(pattern string, f jhttp.ResponseHandler)
	Get(pattern string, f jhttp.ResponseHandler)
	Head(pattern string, f jhttp.ResponseHandler)
	Options(pattern string, f jhttp.ResponseHandler)
	Patch(pattern string, f jhttp.ResponseHandler)
	Post(pattern string, f jhttp.ResponseHandler)
	Put(pattern string, f jhttp.ResponseHandler)
	Trace(pattern string, f jhttp.ResponseHandler)
}

// Mux returns a new Mux that wraps chi.Router functions
type Mux struct {
	Router  chi.Router
	wrapper HTTPWrapper
}

// NewRouter returns a new Mux
func NewRouter(r chi.Router, wrapper HTTPWrapper) Mux {
	return Mux{Router: r, wrapper: wrapper}
}

func (rt Mux) router(r chi.Router) Mux {
	dup := rt
	dup.Router = r
	return dup
}

func (rt Mux) requestWrap(f jhttp.RequestHandler) func(http.Handler) http.Handler {
	return jhttp.RequestWrap(rt.wrapper.RequestWrap(f))
}

func (rt Mux) responseWrap(f jhttp.ResponseHandler) http.HandlerFunc {
	return jhttp.ResponseWrap(rt.wrapper.ResponseWrap(f))
}

func (rt Mux) Use(middlewares ...jhttp.RequestHandler) {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = rt.requestWrap(middleware)
	}
	rt.Router.Use(wrapped...)
}
func (rt Mux) With(middlewares ...jhttp.RequestHandler) Router {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = rt.requestWrap(middleware)
	}
	return rt.router(rt.Router.With(wrapped...))
}

func (rt Mux) Group(fn func(r Router)) Router {
	r := rt.Router.Group(func(r chi.Router) {
		fn(rt.router(r))
	})
	return rt.router(r)
}
func (rt Mux) Route(pattern string, fn func(r Router)) Router {
	r := rt.Router.Route(pattern, func(r chi.Router) {
		fn(rt.router(r))
	})
	return rt.router(r)
}

func (rt Mux) HandleFunc(pattern string, f jhttp.ResponseHandler) {
	rt.Router.HandleFunc(pattern, rt.responseWrap(f))
}
func (rt Mux) MethodFunc(method, pattern string, f jhttp.ResponseHandler) {
	rt.Router.MethodFunc(method, pattern, rt.responseWrap(f))
}

func (rt Mux) Connect(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Connect(pattern, rt.responseWrap(f))
}
func (rt Mux) Delete(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Delete(pattern, rt.responseWrap(f))
}
func (rt Mux) Get(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Get(pattern, rt.responseWrap(f))
}
func (rt Mux) Head(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Head(pattern, rt.responseWrap(f))
}
func (rt Mux) Options(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Options(pattern, rt.responseWrap(f))
}
func (rt Mux) Patch(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Patch(pattern, rt.responseWrap(f))
}
func (rt Mux) Post(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Post(pattern, rt.responseWrap(f))
}
func (rt Mux) Put(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Put(pattern, rt.responseWrap(f))
}
func (rt Mux) Trace(pattern string, f jhttp.ResponseHandler) {
	rt.Router.Trace(pattern, rt.responseWrap(f))
}
