package storage

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
)

const testUUID = "123e4567-e89b-12d3-a456-426614174000"

type uuidTest struct{}

func (u uuidTest) Generate() (string, error) { return testUUID, nil }

type testAPI struct {
	storeMap map[string]string
}

func newTestAPI() testAPI {
	return testAPI{storeMap: make(map[string]string)}
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

var getFilesDir = path.Join(os.TempDir(), test.GenerateName("api_test.Get"))

func (t testAPI) Get(key string) (fs.File, error) {
	content, exists := t.storeMap[key]
	if !exists {
		return nil, fmt.Errorf("key does not exist: %v", key)
	}
	if err := os.MkdirAll(getFilesDir, ioutil.OwnerRWXGroupRX); err != nil {
		return nil, err
	}
	p := filepath.Join(getFilesDir, filepath.Base(key))
	if err := os.WriteFile(p, []byte(content), ioutil.OwnerGroupR); err != nil {
		return nil, err
	}
	return os.Open(p) //nolint:gosec // tests only
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

func (t testAPI) Store(key string, file io.Reader) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	t.storeMap[key] = string(data)
	return nil
}
