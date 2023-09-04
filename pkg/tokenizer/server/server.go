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

	"golang.org/x/exp/slog"
)

// TokenizerServer provides an interface to call tokenizer servers
type TokenizerServer interface {
	Start() error
	Stop() error
	StopAndWait() error
	ForceStop() error
	IsRunning() bool

	Tokenize(str string, obj any) error
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

// CmdOptions is the set of options for the Cmd
type CmdOptions struct {
	name string
	Dir  string
	Args []string
}

// NewCmdOptions returns a new CmdOptions
func NewCmdOptions(name string) CmdOptions {
	return CmdOptions{name: name, Args: []string{}}
}

// Cmd returns the cmd given the options
func (c *CmdOptions) Cmd() (*exec.Cmd, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, c.name, c.Args...) //nolint:gosec //pretty sure not called often
	cmd.Dir = c.Dir
	return cmd, cancel
}

// NewCmdTokenizerServer returns a new CmdServer
func NewCmdTokenizerServer(cmdOpts CmdOptions, port int, stopWarningDuration time.Duration) *CmdTokenizerServer {
	cmd, cancel := cmdOpts.Cmd()
	return &CmdTokenizerServer{
		cmd:                 cmd,
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
			slog.Warn(fmt.Sprintf("Error after waiting for CmdTokenizerServer: %v", err))
		}
		s.isRunning = false
	}()

	for i := 1; i <= 15; i++ {
		time.Sleep(time.Millisecond * 200)
		resp, err := http.Get(s.uriFor(HealthzPath))
		if err != nil {
			continue
		}
		defer resp.Body.Close() //nolint:errcheck // failing is ok
		if resp.StatusCode != http.StatusOK {
			continue
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		if getFirstLine(string(respBytes)) == "ok" {
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
	initialSleep := time.Second
	time.Sleep(initialSleep)
	go func() {
		i := 0
		for {
			select {
			case <-stopped:
				slog.Info("CmdServer stopped")
				return
			default:
				slog.Warn(fmt.Sprintf("CmdServer server is still running after %v",
					time.Duration(i)*s.stopWarningDuration+initialSleep))
			}
			i++
			time.Sleep(s.stopWarningDuration)
		}
	}()
	go func() {
		forceStopDuration := s.stopWarningDuration * 10
		time.Sleep(forceStopDuration)
		if !s.IsRunning() {
			return
		}
		slog.Warn(fmt.Sprintf("Komoran Server still running after %v, calling ForceStop()", forceStopDuration))
		if err2 := s.ForceStop(); err != nil {
			slog.Warn(fmt.Sprintf("error calling ForceStop(): %v", err2))
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

// Tokenize marshalls tokenizes the string into the obj
func (s *CmdTokenizerServer) Tokenize(str string, obj any) error {
	body, err := json.Marshal(&TokenizeRequest{String: str})
	if err != nil {
		return err
	}

	resp, err := http.Post(s.uriFor(TokenizePath), mime.TypeByExtension(".json"), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // failing is ok
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("java.Server.Tokenize() [%v]: %v", resp.StatusCode, string(respBytes))
	}

	if err := json.Unmarshal(respBytes, obj); err != nil {
		return err
	}
	return nil
}

const baseURI = "http://localhost"

func (s *CmdTokenizerServer) uriFor(path string) string {
	return baseURI + ":" + strconv.Itoa(s.port) + path
}
