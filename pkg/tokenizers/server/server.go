// Package server provides a interface for tokenizer servers
package server

import (
	"bytes"
	"context"
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

// TokenizerServer provides an interface to call tokenizer servers
type TokenizerServer interface {
	Start() error
	Stop() error
	StopAndWait() error
	ForceStop() error
	IsRunning() bool

	Tokenize(str string, resp any) error
}

// CmdTokenizerServer is a server that runs a cmd
type CmdTokenizerServer struct {
	cmd       *exec.Cmd
	stdIn     io.WriteCloser
	isRunning bool
	port      int

	stopWarningDuration time.Duration
	cancel              context.CancelFunc
}

// NewCmdTokenizerServer returns a new CmdServer
func NewCmdTokenizerServer(port int, stopWarningDuration time.Duration, name string, args ...string) *CmdTokenizerServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &CmdTokenizerServer{
		cmd:                 exec.CommandContext(ctx, name, args...),
		port:                port,
		stopWarningDuration: stopWarningDuration,
		cancel:              cancel,
	}
}

// HealthzPath is the path to the health check path - GET
const HealthzPath = "/healthz"

// TokenizePath is the path to the tokenize path - POST
const TokenizePath = "/tokenize"

// TokenizeRequest is the format of the request for TokenizePath
type TokenizeRequest struct {
	String string `json:"string"`
}

// Start starts the CmdServer
func (s *CmdTokenizerServer) Start() error {
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
		response, err := http.Get(s.uriFor(HealthzPath))
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
	return errors.New("timed out waiting for " + HealthzPath)
}

func getFirstLine(str string) string {
	newlineIndex := strings.IndexByte(str, '\n')
	if newlineIndex == -1 {
		return str
	}
	return str[:newlineIndex]
}

// Stop stops the CmdServer
func (s *CmdTokenizerServer) Stop() error {
	stopped, err := s.stop()
	go func() {
		i := 0
		for {
			i++
			time.Sleep(s.stopWarningDuration)
			select {
			case <-stopped:
				fmt.Println("CmdServer stopped")
				return
			default:
				fmt.Printf("CmdServer server is still running after %v\n",
					time.Duration(i)*s.stopWarningDuration)
			}
		}
	}()
	go func() {
		forceStopDuration := (s.stopWarningDuration * 10) - time.Second
		time.Sleep(forceStopDuration)
		if !s.IsRunning() {
			return
		}
		fmt.Printf("Komoran Server still running after %v, ForceStop()\n", forceStopDuration)
		if err2 := s.ForceStop(); err != nil {
			fmt.Println(err2)
		}
	}()
	return err
}

// StopKeyword is the keyword to stop the server from stdin
const StopKeyword = "stop"

func (s *CmdTokenizerServer) stop() (chan bool, error) {
	if _, err := io.WriteString(s.stdIn, StopKeyword+"\n"); err != nil {
		return nil, err
	}

	stopped := make(chan bool)
	go func() {
		sleepTime := time.Millisecond * 200
		for i := 1; s.isRunning; i++ {
			time.Sleep(sleepTime)
		}
		stopped <- true
	}()

	return stopped, nil
}

// StopAndWait runs Stop() and waits until the server is stopped
func (s *CmdTokenizerServer) StopAndWait() error {
	_, err := s.stop()
	if err != nil {
		return err
	}

	sleepTime := time.Millisecond * 200
	count := int(s.stopWarningDuration / sleepTime)
	for i := 1; i <= count && s.isRunning; i++ {
		time.Sleep(sleepTime)
	}
	if s.isRunning {
		return fmt.Errorf("CmdServer running after timeout Stop()")
	}
	return nil
}

// ForceStop forces the server to stop via kill
func (s *CmdTokenizerServer) ForceStop() error {
	if s.isRunning {
		return fmt.Errorf("will not ForceStop() while IsRunning()")
	}
	s.cancel()
	return nil
}

// IsRunning returns true if the CmdServer is running
func (s *CmdTokenizerServer) IsRunning() bool {
	return s.isRunning
}

// Tokenize marshalls tokenizes the string into the resp
func (s *CmdTokenizerServer) Tokenize(str string, resp any) error {
	body, err := json.Marshal(&TokenizeRequest{String: str})
	if err != nil {
		return err
	}

	response, err := http.Post(s.uriFor(TokenizePath), mime.TypeByExtension(".json"), bytes.NewBuffer(body))
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

func (s *CmdTokenizerServer) uriFor(path string) string {
	return baseURI + ":" + strconv.Itoa(s.port) + path
}
