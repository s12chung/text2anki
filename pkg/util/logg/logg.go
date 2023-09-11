// Package logg contains utils for slog
package logg

import "log/slog"

// Default returns the default log
func Default() *slog.Logger { return slog.Default() } //nolint:forbidigo // this package abstracts slog

// Err returns an Attr for the error
func Err(err error) slog.Attr { return slog.String("err", err.Error()) }
