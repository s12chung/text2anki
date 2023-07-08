// Package httputil contains utils for http requests
package httputil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/exp/slog"
)

// HTTPError represents an http error
type HTTPError struct {
	Code  int
	Cause error
}

// Error returns the string representation of HTTPError
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Cause)
}

// Error is a safe shorthand to create a new HTTPError
func Error(code int, cause error) *HTTPError {
	return &HTTPError{Code: code, Cause: cause}
}

// ContextKey is a type used for Context Keys
type ContextKey string

// RequestWrapFunc is the function format for RequestWrap, used to set the request (for setting the context)
type RequestWrapFunc func(r *http.Request) (*http.Request, *HTTPError)

// RequestWrap wraps a function that sets the request (for setting the context)
func RequestWrap(f RequestWrapFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newReq, httpError := f(r)
			if httpError != nil {
				RespondError(w, httpError)
				return
			}
			next.ServeHTTP(w, newReq)
		})
	}
}

// RespondJSONWrapFunc is the function format for RespondJSONWrap, used to automatically handle JSON responses
type RespondJSONWrapFunc func(r *http.Request) (any, *HTTPError)

// RespondJSONWrap wraps a function that handles the request using return statements rather than writing to the response
func RespondJSONWrap(f RespondJSONWrapFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, httpError := f(r)
		if httpError != nil {
			RespondError(w, httpError)
			slog.Error(httpError.Cause.Error(), slog.Int("code", httpError.Code), slog.String("url", r.URL.Path), slog.String("method", r.Method))
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
func RespondError(w http.ResponseWriter, httpError *HTTPError) {
	w.Header().Set("Content-Type", JSONContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpError.Code)

	err := json.NewEncoder(w).Encode(&ErrResponse{
		Error:      httpError.Cause.Error(),
		Code:       httpError.Code,
		StatusText: http.StatusText(httpError.Code),
	})
	if err != nil {
		http.Error(w, httpError.Error(), httpError.Code)
	}
}

// RespondJSON responds with JSON
func RespondJSON(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", JSONContentType)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		RespondError(w, Error(http.StatusInternalServerError, err))
	}
}

// ExtractJSON binds the JSON request body to the given struct
func ExtractJSON(r *http.Request, o any) *HTTPError {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Error(http.StatusBadRequest, err)
	}
	err = json.Unmarshal(body, o)
	if err != nil {
		return Error(http.StatusBadRequest, err)
	}
	return nil
}

// ReturnModelOr500 runs the modelFunc, and returns http.StatusInternalServerError for the error
func ReturnModelOr500(modelFunc func() (any, error)) (any, *HTTPError) {
	model, err := modelFunc()
	if err != nil {
		return nil, Error(http.StatusInternalServerError, err)
	}
	return model, nil
}
