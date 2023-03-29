package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/s12chung/text2anki/integrations/tokenizers/khaiii/pkg/khaiii"
	"github.com/s12chung/text2anki/pkg/tokenizers/server/serverimpl"
)

var port int

func init() {
	serverimpl.FsPort(&port, flag.CommandLine)
	flag.Parse()
}

func main() {
	if err := run(port); err != nil {
		fmt.Println(err)
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
	server := serverimpl.NewServerImpl(khaiii.NewTokenizer(k))
	return server.Run(port)
}
