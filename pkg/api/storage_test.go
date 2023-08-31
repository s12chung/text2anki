package api

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/test"
)

var storageServer test.Server

func init() {
	storageServer = server.WithPathPrefix(config.StorageURLPath)
}

func TestRoutes_StoragePut(t *testing.T) {
	testName := "TestRoutes_StoragePut"
	req, err := routes.Storage.DBStorage.SignPut("test_table", "my_column", ".txt")
	require.NoError(t, err)

	u, err := url.Parse(req.URL)
	require.NoError(t, err)
	key := strings.TrimPrefix(u.Path, config.StorageURLPath)

	badQuery := localstore.CipherQueryParam + "=" + base64.URLEncoding.EncodeToString([]byte("my_bad"))

	testCases := []struct {
		name         string
		path         string
		body         []byte
		expectedCode int
		key          string
	}{
		{name: "normal", path: key + "?" + u.RawQuery, key: key, body: []byte("test me"), expectedCode: http.StatusOK},
		{name: "bad_query", path: "/testy?" + badQuery, body: []byte("test me"), expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, storageServer.NewRequest(t, http.MethodPut, tc.path, bytes.NewReader(tc.body)))
			resp.EqualCode(t, tc.expectedCode)
			testModelResponse(t, resp, testName, tc.name, &StoragePutOk{})

			if resp.Code != http.StatusOK {
				return
			}
			fileBytes, err := os.ReadFile(path.Join(routesConfig.StorageConfig.LocalStoreConfig.KeyBasePath, tc.key))
			require.NoError(err)
			require.Equal(string(tc.body), string(fileBytes))
		})
	}
}

func TestRoutes_StorageGet(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		body         []byte
		expectedCode int
	}{
		{name: "normal", key: "/my_table/here/go.txt", body: []byte("test me"), expectedCode: http.StatusOK},
		{name: "not_exist", key: "/another/dir/hi.txt", body: nil, expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			if tc.body != nil {
				require.NoError(routes.Storage.Storer.Store(tc.key, bytes.NewReader(tc.body)))
			}
			resp := test.HTTPDo(t, storageServer.NewRequest(t, http.MethodGet, tc.key, bytes.NewReader(tc.body)))
			resp.EqualCode(t, tc.expectedCode)
			if resp.Code != http.StatusOK {
				return
			}
			require.Equal(resp.Body.String(), string(tc.body))
		})
	}
}
