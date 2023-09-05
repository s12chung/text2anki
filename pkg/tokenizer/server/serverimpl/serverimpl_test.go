package serverimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/tokenizer/server"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/logg"
	"github.com/s12chung/text2anki/pkg/util/test"
)

const host = "http://localhost"
const testPort = 9000

func TestMain(m *testing.M) {
	server := NewServerImpl(&SplitTokenizer{})
	serverChannel := server.runWithoutStdin(testPort)
	go func() {
		err := <-serverChannel
		if err != nil {
			slog.Error("serverimpl serverChannel", logg.Err(err))
			os.Exit(-1)
		}
	}()
	code := m.Run()
	if err := server.Stop(); err != nil {
		slog.Error("serverimpl server.Stop()", logg.Err(err))
	}
	if !cleaned {
		slog.Error("cleaned = false from Cleanup()")
		os.Exit(-1)
	}
	os.Exit(code)
}

var cleaned = false

type SplitTokenizer struct{}

func (s *SplitTokenizer) Cleanup() { cleaned = true }
func (s *SplitTokenizer) Tokenize(str string) (any, error) {
	return &tokenizeResponse{strings.Split(str, " ")}, nil
}

type tokenizeResponse struct {
	Tokens []string `json:"tokens"`
}

func TestHealthz(t *testing.T) {
	require := require.New(t)

	resp, err := httputil.Get(context.Background(), getURI(server.HealthzPath))
	require.NoError(err)
	defer func() { require.NoError(resp.Body.Close()) }()

	require.Equal(http.StatusOK, resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")
	require.Equal("text/plain; charset=utf-8", contentType)

	data, err := io.ReadAll(resp.Body)
	require.NoError(err)

	require.True(strings.HasPrefix(string(data), "ok\n"))
}

func TestTokenize(t *testing.T) {
	require := require.New(t)
	input := server.TokenizeRequest{
		String: "my example",
	}

	resp, err := httputil.Post(context.Background(), getURI(server.TokenizePath),
		jhttp.JSONContentType,
		bytes.NewBuffer(test.JSON(t, input)))
	require.NoError(err)
	defer func() { require.NoError(resp.Body.Close()) }()

	require.Equal(http.StatusOK, resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")
	require.Equal(jhttp.JSONContentType, contentType)

	data := &tokenizeResponse{}
	err = json.NewDecoder(resp.Body).Decode(data)
	require.NoError(err)

	expectedTokens := []string{"my", "example"}
	require.Equal(expectedTokens, data.Tokens)
}

func getURI(path string) string {
	return fmt.Sprintf(host+":%v%v", testPort, path)
}
