// Package httputilchi provides a light chi.Router wrapper that wraps chi.Router functions and adheres to
// httputil.RespondJSONWrapFunc and httputil.RequestWrapFunc
// nolint: revive
package httputilchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/s12chung/text2anki/pkg/util/httputil"
)

// HTTPWrapper provides wrapping functions that Router uses
type HTTPWrapper interface {
	WrapRequest(f httputil.RequestWrapFunc) httputil.RequestWrapFunc
	WrapResponse(f httputil.RespondJSONWrapFunc) httputil.RespondJSONWrapFunc
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

func (t Router) wrapRequest(f httputil.RequestWrapFunc) func(http.Handler) http.Handler {
	return httputil.RequestWrap(t.wrapper.WrapRequest(f))
}

func (t Router) wrapResponse(f httputil.RespondJSONWrapFunc) http.HandlerFunc {
	return httputil.RespondJSONWrap(t.wrapper.WrapResponse(f))
}

func (t Router) Use(middlewares ...httputil.RequestWrapFunc) {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = t.wrapRequest(middleware)
	}
	t.Router.Use(wrapped...)
}
func (t Router) With(middlewares ...httputil.RequestWrapFunc) Router {
	wrapped := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		wrapped[i] = t.wrapRequest(middleware)
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

func (t Router) HandleFunc(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.HandleFunc(pattern, t.wrapResponse(f))
}
func (t Router) MethodFunc(method, pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.MethodFunc(method, pattern, t.wrapResponse(f))
}

func (t Router) Connect(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Connect(pattern, t.wrapResponse(f))
}
func (t Router) Delete(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Delete(pattern, t.wrapResponse(f))
}
func (t Router) Get(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Get(pattern, t.wrapResponse(f))
}
func (t Router) Head(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Head(pattern, t.wrapResponse(f))
}
func (t Router) Options(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Options(pattern, t.wrapResponse(f))
}
func (t Router) Patch(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Patch(pattern, t.wrapResponse(f))
}
func (t Router) Post(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Post(pattern, t.wrapResponse(f))
}
func (t Router) Put(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Put(pattern, t.wrapResponse(f))
}
func (t Router) Trace(pattern string, f httputil.RespondJSONWrapFunc) {
	t.Router.Trace(pattern, t.wrapResponse(f))
}
