// Package serverimpl provides a wrapper implementation fo a Tokenizer server
package serverimpl

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/s12chung/text2anki/pkg/tokenizers/server"
)

// Tokenizer is the interface for the Tokenizer that the Server works with
type Tokenizer interface {
	Cleanup()
	Tokenize(str string) (any, error)
}

// ServerImpl is an implemenation of a Tokenizer server
type ServerImpl struct {
	tokenizer Tokenizer
	server    *http.Server
	listener  net.Listener

	running bool
}

// NewServerImpl returns a new ServerImpl
func NewServerImpl(tokenizer Tokenizer) *ServerImpl {
	return &ServerImpl{
		tokenizer: tokenizer,
	}
}

const defaultPort = 9999

// FsPort adds a flag for the port in the FlagSet
func FsPort(port *int, fs *flag.FlagSet) {
	fs.IntVar(port, "port", defaultPort, "the port number to listen on")
}

// Run runs the server
func (s *ServerImpl) Run(port int) error {
	serverChannel := s.runWithoutStdin(port)
	return s.waitStdinStop(serverChannel)
}

func (s *ServerImpl) runWithoutStdin(port int) chan error {
	serverChannel := make(chan error)
	if err := s.setupServer(port); err != nil {
		serverChannel <- err
	}
	go func() {
		s.running = true
		serverChannel <- s.server.Serve(s.listener)
	}()

	retChannel := make(chan error)
	go func() {
		err := <-serverChannel
		if err == http.ErrServerClosed {
			err = nil
		}
		s.running = false
		retChannel <- err
	}()
	return retChannel
}

// Stop stops the server
func (s *ServerImpl) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *ServerImpl) setupServer(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc(server.HealthzPath, handleHeathzfunc)
	mux.HandleFunc(server.TokenizePath, handleWrapper(s.handleTokenize))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}
	server.RegisterOnShutdown(s.tokenizer.Cleanup)

	addr := server.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.server, s.listener = server, ln
	return nil
}

func (s *ServerImpl) waitStdinStop(serverChannel chan error) error {
	manualStop := false

	var err2 error
	go func() {
		err2 = <-serverChannel
		if manualStop {
			return
		}

		_, err3 := io.WriteString(os.Stdin, server.StopKeyword+"\n")
		if err2 != nil {
			err2 = err3
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() != server.StopKeyword {
			continue
		}
		break
	}
	manualStop = true
	if err := s.Stop(); err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return scanner.Err()
}

func handleHeathzfunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".txt"))
	fmt.Fprintf(w, "ok\n%s", time.Now().Format(time.RFC3339))
}

type handleWrapperFunc func(r *http.Request) ([]byte, int, error)

func handleWrapper(f handleWrapperFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, status, err := f(r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		w.Header().Set("Content-Type", mime.TypeByExtension(".json"))
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type tokenizeResponse struct {
	Tokens any `json:"tokens"`
}

func (s *ServerImpl) handleTokenize(r *http.Request) ([]byte, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("405 Method Not Allowed")
	}

	req := &server.TokenizeRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	tokens, err := s.tokenizer.Tokenize(req.String)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	b, err := json.Marshal(tokenizeResponse{Tokens: tokens})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return b, 0, nil
}