// Package httputil contains utils for http requests
package httputil

import (
	"encoding/json"
	"io"
	"net/http"
)

// ContextKey is a type used for Context Keys
type ContextKey string

// RequestWrapFunc is the function format for RequestWrap, used to set the request (for setting the context)
type RequestWrapFunc func(r *http.Request) (*http.Request, int, error)

// RequestWrap wraps a function that sets the request (for setting the context)
func RequestWrap(f RequestWrapFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newReq, code, err := f(r)
			if err != nil {
				RespondError(w, code, err)
				return
			}
			next.ServeHTTP(w, newReq)
		})
	}
}

// RespondJSONWrapFunc is the function format for RespondJSONWrap, used to automatically handle JSON responses
type RespondJSONWrapFunc func(r *http.Request) (any, int, error)

// RespondJSONWrap wraps a function that handles the request using return statements rather than writing to the response
func RespondJSONWrap(f RespondJSONWrapFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := f(r)
		if err != nil {
			RespondError(w, code, err)
			return
		}
		RespondJSON(w, resp)
	}
}

// ErrResponse is the struct used for the JSON error response
type ErrResponse struct {
	Error      string `json:"error,omitempty"`
	Code       int    `json:"code,omitempty"`
	StatusText string `json:"status_text,omitempty"`
}

// JSONContentType is the content type used for the JSON responses
const JSONContentType = "application/json; charset=utf-8"

// RespondError responds with a JSON error
func RespondError(w http.ResponseWriter, code int, sourceError error) {
	w.Header().Set("Content-Type", JSONContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(&ErrResponse{
		Error:      sourceError.Error(),
		Code:       code,
		StatusText: http.StatusText(code),
	})
	if err != nil {
		http.Error(w, sourceError.Error(), code)
	}
}

// RespondJSON responds with JSON
func RespondJSON(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", JSONContentType)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		RespondError(w, http.StatusInternalServerError, err)
	}
}

// BindJSON binds the JSON request body to the given struct
func BindJSON(r *http.Request, o any) (int, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = json.Unmarshal(body, o)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}
