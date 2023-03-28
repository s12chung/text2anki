package serverimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/s12chung/text2anki/pkg/tokenizers/server"
	"github.com/stretchr/testify/require"
)

type TokenResponse struct {
	Tokens []string `json:"tokens"`
}

type SplitTokenizer struct {
}

var cleaned = false

func (s *SplitTokenizer) Cleanup() {
	cleaned = true
}
func (s *SplitTokenizer) Tokenize(str string) (any, error) {
	return strings.Split(str, " "), nil
}

const host = "http://localhost"

func TestMain(m *testing.M) {
	server := NewServerImpl(&SplitTokenizer{})
	serverChannel := server.runWithoutStdin(defaultPort)
	go func() {
		err := <-serverChannel
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}()
	code := m.Run()
	if err := server.Stop(); err != nil {
		fmt.Println(err)
	}
	if !cleaned {
		fmt.Println("cleaned = false from Cleanup()")
		os.Exit(-1)
	}
	os.Exit(code)
}

func TestHealthz(t *testing.T) {
	require := require.New(t)

	resp, err := http.Get(fmt.Sprintf(host+":%v%v", defaultPort, server.HealthzPath))
	require.NoError(err)
	defer func() {
		require.NoError(resp.Body.Close())
	}()

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
	payload, err := json.Marshal(input)
	require.NoError(err)

	resp, err := http.Post(fmt.Sprintf(host+":%v%v", defaultPort, server.TokenizePath),
		mime.TypeByExtension(".json"),
		bytes.NewBuffer(payload))
	require.NoError(err)
	defer func() {
		require.NoError(resp.Body.Close())
	}()

	require.Equal(http.StatusOK, resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")
	require.Equal("application/json", contentType)

	data := &TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(data)
	require.NoError(err)

	expectedTokens := []string{"my", "example"}
	require.Equal(expectedTokens, data.Tokens)
}
