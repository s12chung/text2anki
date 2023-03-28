// Package java includes helpers to execute java server
package java

import (
	"strconv"
	"time"

	"github.com/s12chung/text2anki/pkg/tokenizers/server"
)

// NewJarServer returns a server that runs a jar file
func NewJarServer(dir, jarName string, port, backlog int, stopWarningDuration time.Duration) server.TokenizerServer {
	return server.NewCmdTokenizerServer(port, stopWarningDuration,
		dir,
		"java", "-jar", jarName,
		strconv.Itoa(port), strconv.Itoa(backlog))
}
