package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// JSON returns indented json
func JSON(t *testing.T, v any) []byte {
	require := require.New(t)
	bytes, err := json.MarshalIndent(v, "", "  ")
	require.NoError(err)
	return bytes
}

// IndentJSON returns indented JSON
func IndentJSON(t *testing.T, original []byte) []byte {
	require := require.New(t)
	var changed bytes.Buffer
	err := json.Indent(&changed, original, "", "  ")
	require.NoError(err)
	return changed.Bytes()
}

// Read reads the file
func Read(t *testing.T, filepath string) []byte {
	require := require.New(t)
	expected, err := os.ReadFile(filepath) //nolint:gosec // for tests
	require.NoError(err)
	return []byte(strings.TrimSpace(string(expected)))
}

// ReadBytes reads bytes from readCloser
func ReadBytes(t *testing.T, readCloser io.ReadCloser) []byte {
	require := require.New(t)
	b, err := io.ReadAll(readCloser)
	require.NoError(err)
	return b
}

// Unmarshall is a one-liner to json unmarshall
func Unmarshall(t *testing.T, b []byte, data any) {
	require := require.New(t)
	require.NoError(json.Unmarshal(b, data))
}

// StaticCopyable returns a TestCopy of a struct
type StaticCopyable interface {
	StaticCopy() any
}

// StaticCopy returns a static JSON copy
func StaticCopy(t *testing.T, b []byte, data StaticCopyable) []byte {
	Unmarshall(t, b, data)
	return JSON(t, data.StaticCopy())
}

// StaticCopyOrIndent returns a static JSON copy when http.StatusOK or indented JSON copy
func StaticCopyOrIndent(t *testing.T, code int, b []byte, data StaticCopyable) []byte {
	if code == http.StatusOK {
		return StaticCopy(t, b, data)
	}
	return IndentJSON(t, b)
}

// HTTPDo does a http.DefaultClient.Do and returns a Response
func HTTPDo(t *testing.T, req *http.Request) Response {
	require := require.New(t)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(err)
	return NewResponse(t, resp)
}

// Response is a wrapper around http.Response, similar to httptest.ResponseRecorder
type Response struct {
	Code   int
	Body   *bytes.Buffer
	result *http.Response
}

// NewResponse returns a new Response
func NewResponse(t *testing.T, resp *http.Response) Response {
	return Response{
		Code:   resp.StatusCode,
		Body:   bytes.NewBuffer(ReadBytes(t, resp.Body)),
		result: resp,
	}
}

// Result returns the internal http.Response
func (r Response) Result() *http.Response {
	return r.result
}
