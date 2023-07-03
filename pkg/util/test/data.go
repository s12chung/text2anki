package test

import (
	"bytes"
	"encoding/json"
	"io"
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

// StaticCopySlice returns a static JSON copy of a datas slice
func StaticCopySlice(t *testing.T, b []byte, datas any) []byte {
	require := require.New(t)

	value := reflect.ValueOf(datas)
	typ := value.Type()

	require.True(typ.Kind() == reflect.Ptr)
	require.True(typ.Elem().Kind() == reflect.Slice)
	require.True(typ.Elem().Elem().Implements(reflect.TypeOf((*StaticCopyable)(nil)).Elem()))

	Unmarshall(t, b, datas)

	sliceValue := value.Elem()
	length := sliceValue.Len()
	staticCopies := make([]any, length)
	for i := 0; i < length; i++ {
		element := sliceValue.Index(i)
		copyable, ok := element.Interface().(StaticCopyable)
		if !ok {
			require.Fail("Element is not StaticCopyable (should never happen due to check above)")
		}
		staticCopies[i] = copyable.StaticCopy()
	}
	return JSON(t, staticCopies)
}

// StaticCopyOrIndentSlice returns a static JSON copy when http.StatusOK or indented JSON copy
func StaticCopyOrIndentSlice(t *testing.T, code int, b []byte, datas any) []byte {
	if code == http.StatusOK {
		return StaticCopySlice(t, b, datas)
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
	for i := 0; i < v.NumField(); i++ {
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

// Server is wrapper around httptest.Server for convenience
type Server struct {
	pathPrefix string
	*httptest.Server
}

// NewRequest returns a new request for the server
func (s Server) NewRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	require := require.New(t)
	req, err := http.NewRequest(method, s.URL+s.pathPrefix+path, body)
	require.NoError(err)
	return req
}

// WithPathPrefix returns a new server with the pathPrefix set for NewRequest
func (s Server) WithPathPrefix(prefix string) Server {
	dup := s
	dup.pathPrefix += prefix
	return dup
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
