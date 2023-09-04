// Package httputilchi provides a light chi.Router wrapper that wraps chi.Router functions and adheres to
// httputil.ResponseJSONWrapFunc and httputil.RequestWrapFunc
// nolint: revive
package httputilchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/util/httputil"
)

// HTTPWrapper provides wrapping functions that Router uses
type HTTPWrapper interface {
	RequestWrap(f httputil.RequestWrapFunc) httputil.RequestWrapFunc
	ResponseWrap(f httputil.ResponseJSONWrapFunc) httputil.ResponseJSONWrapFunc
}

// Router returns a new Router that wraps chi.Router functions
type Router struct {
	Router  chi.Router
	wrapper HTTPWrapper
}

// NewRouter returns a new Router
func NewRouter(r chi.Router, wrapper HTTPWrapper) Router {
	return Router{Router: r, wrapper: wrapper}
}

func (t Router) router(r chi.Router) Router {
	dup := t
	dup.Router = r
	return dup
}

func (t Router) requestWrap(f httputil.RequestWrapFunc) func(http.Handler) http.Handler {
	return httputil.RequestWrap(t.wrapper.RequestWrap(f))
}

func (t Router) responseWrap(f httputil.ResponseJSONWrapFunc) http.HandlerFunc {
	return httputil.ResponseJSONWrap(t.wrapper.ResponseWrap(f))
}

func (t Router) Use(middlewares ...httputil.RequestWrapFunc) {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = t.requestWrap(middleware)
	}
	t.Router.Use(wrapped...)
}
func (t Router) With(middlewares ...httputil.RequestWrapFunc) Router {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = t.requestWrap(middleware)
	}
	return t.router(t.Router.With(wrapped...))
}

func (t Router) Group(fn func(r Router)) Router {
	r := t.Router.Group(func(r chi.Router) {
		fn(t.router(r))
	})
	return t.router(r)
}
func (t Router) Route(pattern string, fn func(r Router)) Router {
	r := t.Router.Route(pattern, func(r chi.Router) {
		fn(t.router(r))
	})
	return t.router(r)
}

func (t Router) HandleFunc(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.HandleFunc(pattern, t.responseWrap(f))
}
func (t Router) MethodFunc(method, pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.MethodFunc(method, pattern, t.responseWrap(f))
}

func (t Router) Connect(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Connect(pattern, t.responseWrap(f))
}
func (t Router) Delete(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Delete(pattern, t.responseWrap(f))
}
func (t Router) Get(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Get(pattern, t.responseWrap(f))
}
func (t Router) Head(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Head(pattern, t.responseWrap(f))
}
func (t Router) Options(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Options(pattern, t.responseWrap(f))
}
func (t Router) Patch(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Patch(pattern, t.responseWrap(f))
}
func (t Router) Post(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Post(pattern, t.responseWrap(f))
}
func (t Router) Put(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Put(pattern, t.responseWrap(f))
}
func (t Router) Trace(pattern string, f httputil.ResponseJSONWrapFunc) {
	t.Router.Trace(pattern, t.responseWrap(f))
}
