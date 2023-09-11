// Package flog has slog helpers with fixture
package flog

import (
	"log/slog"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

type noWriter struct{}

func (n noWriter) Write(p []byte) (int, error) { return len(p), nil }

// FixtureUpdateNoWrite returns a logger that does nothing when fixture.WillUpdate() is true
func FixtureUpdateNoWrite() *slog.Logger {
	if fixture.WillUpdate() {
		return slog.New(slog.NewTextHandler(noWriter{}, nil))
	}
	return slog.Default() //nolint:forbidigo // this package abstracts slog
}
