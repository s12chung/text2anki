// Package logg contains utils for slog
package logg

import "log/slog"

// Err returns an Attr for the error
func Err(err error) slog.Attr { return slog.String("err", err.Error()) }
