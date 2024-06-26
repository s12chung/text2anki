// Package jhttp reformats http handlers to follow returning error conventions by abstracting to JSON responses
package jhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/s12chung/text2anki/pkg/util/logg"
)

var plog *slog.Logger

// SetLog setts the log for the package
func SetLog(log *slog.Logger) { plog = log }

// HTTPError represents an http error
type HTTPError struct {
	Code  int
	Cause error
}

// Error returns the string representation of HTTPError
func (e HTTPError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Cause)
}

// Error is a safe shorthand to create a new HTTPError
func Error(code int, cause error) *HTTPError {
	return &HTTPError{Code: code, Cause: cause}
}

// ContextKey is a type used for Context Keys
type ContextKey string

// RequestHandler is the function format for RequestWrap, used to set the request (for setting the context)
type RequestHandler func(r *http.Request) (*http.Request, *HTTPError)

// RequestWrap wraps a function that sets the request (for setting the context)
func RequestWrap(f RequestHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, httpError := f(r)
			if httpError != nil {
				RespondError(w, httpError)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

// ResponseHandler is the function format for ResponseWrap, used to automatically handle JSON responses
type ResponseHandler func(r *http.Request) (any, *HTTPError)

func (res ResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ResponseWrap(res)(w, r)
}

// ResponseWrap wraps a function that handles the request using return statements rather than writing to the response
func ResponseWrap(f ResponseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model, httpErr := f(r)
		if httpErr != nil {
			LogAndRespondError(w, r, httpErr)
			return
		}
		RespondJSON(w, model)
	}
}

// LogAndRespondError logs and responds the error
func LogAndRespondError(w http.ResponseWriter, r *http.Request, httpErr *HTTPError) {
	LogError(r, httpErr)
	RespondError(w, httpErr)
}

// LogError logs the error
func LogError(r *http.Request, httpErr *HTTPError) {
	plog.LogAttrs(r.Context(), slog.LevelError, "jhttp response", append(logg.RequestAttrs(r), slog.Int("code", httpErr.Code), logg.Err(httpErr))...)
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
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, Error(http.StatusInternalServerError, err)
	}
	return model, nil
}

// ReturnSliceOr500 runs the sliceFunc, and returns http.StatusInternalServerError for the error
func ReturnSliceOr500[T any](sliceFunc func() ([]T, error)) (any, *HTTPError) {
	return ReturnModelOr500(func() (any, error) {
		slice, err := sliceFunc()
		if slice == nil {
			return make([]T, 0), err
		}
		return slice, err
	})
}
