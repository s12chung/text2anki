// Package server provides a interface for tokenizer servers
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Server provides an interface to call tokenizer servers
type Server interface {
	Start() error
	Stop() error
	IsRunning() bool

	Tokenize(str string, resp any) error
}

// CmdServer is a server that runs a cmd
type CmdServer struct {
	cmd       *exec.Cmd
	stdIn     io.WriteCloser
	isRunning bool
	port      int
}

// NewCmdSever returns a new CmdServer
func NewCmdSever(port int, name string, args ...string) *CmdServer {
	return &CmdServer{
		cmd:  exec.Command(name, args...),
		port: port,
	}
}

const healthzPath = "/healthz"

// Start starts the CmdServer
func (s *CmdServer) Start() error {
	var err error

	s.stdIn, err = s.cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = s.cmd.Start()
	if err != nil {
		return err
	}
	s.isRunning = true

	go func() {
		if err := s.cmd.Wait(); err != nil {
			fmt.Printf("Error after waiting for command: %v\n", err)
		}
		s.isRunning = false
	}()

	for i := 1; i <= 15; i++ {
		time.Sleep(time.Millisecond * 200)
		response, err := http.Get(s.uriFor(healthzPath))
		if err != nil {
			continue
		}
		if response.StatusCode != http.StatusOK {
			continue
		}
		resp, err := io.ReadAll(response.Body)
		if err != nil {
			continue
		}
		if getFirstLine(string(resp)) == "ok" {
			return nil
		}
	}
	return errors.New("timed out waiting for " + healthzPath)
}

func getFirstLine(str string) string {
	newlineIndex := strings.IndexByte(str, '\n')
	if newlineIndex == -1 {
		return str
	}
	return str[:newlineIndex]
}

// Stop stops the CmdServer
func (s *CmdServer) Stop() error {
	_, err := io.WriteString(s.stdIn, "stop\n")
	return err
}

// IsRunning returns true if the CmdServer is running
func (s *CmdServer) IsRunning() bool {
	return s.isRunning
}

type tokenizeRequest struct {
	String string `json:"string"`
}

// Tokenize marshalls tokenizes the string into the resp
func (s *CmdServer) Tokenize(str string, resp any) error {
	body, err := json.Marshal(&tokenizeRequest{String: str})
	if err != nil {
		return err
	}

	response, err := http.Post(s.uriFor("/tokenize"), mime.TypeByExtension(".json"), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("java.Server.Tokenize() [%v]: %v", response.StatusCode, string(respBytes))
	}

	if err := json.Unmarshal(respBytes, resp); err != nil {
		return err
	}
	return nil
}

const baseURI = "http://localhost"

func (s *CmdServer) uriFor(path string) string {
	return baseURI + ":" + strconv.Itoa(s.port) + path
}
