// Package logg contains utils for slog
package logg

import (
	"log/slog"
	"net/http"
)

// Default returns the default log
func Default() *slog.Logger { return slog.Default() } //nolint:forbidigo // this package abstracts slog

// Err returns an Attr for the error
func Err(err error) slog.Attr { return slog.String("err", err.Error()) }

// RequestAttrs returns []slog.Attr for the request
func RequestAttrs(r *http.Request) []slog.Attr {
	return []slog.Attr{slog.String("method", r.Method), slog.String("url", r.URL.Path)}
}
