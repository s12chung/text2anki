package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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

// JSONString returns indented json as a string
func JSONString(t *testing.T, v any) string {
	return string(JSON(t, v))
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
type StaticCopyable[T any] interface {
	StaticCopy() T
}

// StaticCopy returns a static JSON copy
func StaticCopy[T StaticCopyable[T]](t *testing.T, b []byte, model *T) []byte {
	Unmarshall(t, b, model)
	return JSON(t, (*model).StaticCopy())
}

// StaticCopyOrIndent returns a static JSON copy when http.StatusOK or indented JSON copy
func StaticCopyOrIndent[T StaticCopyable[T]](t *testing.T, code int, b []byte, model *T) []byte {
	if code == http.StatusOK {
		return StaticCopy[T](t, b, model)
	}
	return IndentJSON(t, b)
}

// StaticCopySlice returns a static JSON copy of a datas slice
func StaticCopySlice[T StaticCopyable[T]](t *testing.T, b []byte, models *[]T) []byte {
	if models == nil {
		models = &[]T{}
	}
	Unmarshall(t, b, models)

	staticCopies := make([]any, len(*models))
	for i, model := range *models {
		staticCopies[i] = model.StaticCopy()
	}
	return JSON(t, staticCopies)
}

// StaticCopyOrIndentSlice returns a static JSON copy when http.StatusOK or indented JSON copy
func StaticCopyOrIndentSlice[T StaticCopyable[T]](t *testing.T, code int, b []byte, models *[]T) []byte {
	if code == http.StatusOK {
		return StaticCopySlice[T](t, b, models)
	}
	return IndentJSON(t, b)
}

// EmptyFields returns a slice of the empty fields of s
func EmptyFields(t *testing.T, s interface{}) []string {
	require := require.New(t)

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	require.Equal(reflect.Struct, v.Kind())

	emptyFields := make([]string, v.NumField())
	a := 0
	for i := range v.NumField() {
		if v.Field(i).IsZero() {
			emptyFields[a] = v.Type().Field(i).Name
			a++
		}
	}
	if a == 0 {
		return nil
	}
	return emptyFields[:a]
}

// EmptyFieldsMatch checks if s has the given empty fields
func EmptyFieldsMatch(t *testing.T, s interface{}, emptyFields ...string) {
	require := require.New(t)
	require.ElementsMatch(emptyFields, EmptyFields(t, s))
}

// Server is wrapper around httptest.Server for convenience
type Server struct {
	pathPrefix string
	*httptest.Server
}

// NewRequest returns a new request for the server
// nolint: revive // prefer testing.T to be first
func (s Server) NewRequest(t *testing.T, ctx context.Context, method, path string, body io.Reader) *http.Request {
	require := require.New(t)
	require.NotNil(s.Server, "Server is not set (due to init timing?)")
	req, err := http.NewRequestWithContext(ctx, method, s.URL+s.pathPrefix+path, body)
	require.NoError(err)
	return req
}

// WithPathPrefix returns a new server with the pathPrefix set for NewRequest
func (s Server) WithPathPrefix(prefix string, log *slog.Logger) Server {
	if s.Server == nil {
		log.Error(fmt.Sprintf("test.Server is not set before calling WithPathPrefix(%v) - due to init timing?", prefix))
		os.Exit(-1)
	}

	dup := s
	dup.pathPrefix += prefix
	return dup
}

// HTTPDo does a http.DefaultClient.Do and returns a Response
func HTTPDo(t *testing.T, req *http.Request) Response {
	require := require.New(t)
	resp, err := http.DefaultClient.Do(req)
	defer func() { require.NoError(resp.Body.Close()) }()
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

// EqualCode tests if the given status code matches the response
func (r Response) EqualCode(t *testing.T, expected int) {
	require.Equal(t, expected, r.Code, "body: "+r.Body.String())
}
