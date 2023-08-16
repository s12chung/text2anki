package storage

import (
	"net/http"
	"path"
	"strings"
)

const testUUID = "123e4567-e89b-12d3-a456-426614174000"

type uuidTest struct {
}

func (u uuidTest) Generate() (string, error) {
	return testUUID, nil
}

type testAPI struct {
}

const keyURLPrefix = "http://localhost:3000/"

func keyURL(key string) string {
	return keyURLPrefix + key
}

func (t testAPI) SignPut(key string) (PreSignedHTTPRequest, error) {
	return PreSignedHTTPRequest{
		URL:          keyURL(key) + "?cipher=blah",
		Method:       "PUT",
		SignedHeader: http.Header{},
	}, nil
}

func (t testAPI) SignGet(key string) (string, error) {
	return keyURL(key), nil
}

func (t testAPI) ListKeys(prefix string) ([]string, error) {
	if prefix != "sources/parts/123e4567-e89b-12d3-a456-426614174000" {
		return []string{}, nil
	}
	return []string{
			path.Join(prefix, "parts.PreParts[0].Image.jpg"),
			path.Join(prefix, "parts.PreParts[0].Audio.mp3"),
			path.Join(prefix, "parts.PreParts[1].Image.png"),
			path.Join(prefix, "parts.PreParts[2].Audio.mp3"),
		},
		nil
}

func (t testAPI) KeyFromSignGet(key string) (string, error) {
	return strings.TrimPrefix(key, keyURLPrefix), nil
}
