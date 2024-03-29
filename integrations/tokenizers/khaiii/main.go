// Package main is the start point for khaiii
package main

import (
	"flag"
	"os"

	"github.com/s12chung/text2anki/integrations/tokenizers/khaiii/pkg/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizer/server/serverimpl"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

var port int
var plog = logg.Default() //nolint:forbidigo // main package

func init() {
	serverimpl.FsPort(&port, flag.CommandLine)
	flag.Parse()
}

func main() {
	if err := run(port); err != nil {
		plog.Error("khaiii/main", logg.Err(err))
		os.Exit(-1)
	}
}

func run(port int) error {
	var err error
	k, err := khaiii.NewKhaiii(khaiii.DefaultDlPath)
	if err != nil {
		return err
	}
	if err = k.Open(khaiii.DefaultRscPath); err != nil {
		return err
	}
	server := serverimpl.NewServerImpl(khaiii.NewTokenizer(k, plog))
	return server.Run(port)
}
