package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/s12chung/text2anki/pkg/tokenizers/server/serverimpl"
	"github.com/s12chung/text2anki/tokenizers/khaiii/pkg/khaiii"
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

type KhaiiiServer struct {
	kahiii *khaiii.Khaiii
}

func NewKhaiiiServer(k *khaiii.Khaiii) *KhaiiiServer {
	return &KhaiiiServer{
		kahiii: k,
	}
}

func (k *KhaiiiServer) Cleanup() {
	if err := k.kahiii.Close(); err != nil {
		fmt.Println(err)
	}
}

func (k *KhaiiiServer) Tokenize(str string) (any, error) {
	return k.kahiii.Analyze(str)
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
	server := serverimpl.NewServerImpl(NewKhaiiiServer(k))
	return server.Run(port)
}
