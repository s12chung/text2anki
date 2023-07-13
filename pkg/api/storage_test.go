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

	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/test"
)

func TestRoutes_StoragePut(t *testing.T) {
	testName := "TestRoutes_StoragePut"
	reqs, err := routes.Storage.Signer.Sign("test_table", "my_field", []string{".txt"})
	require.NoError(t, err)
	u, err := url.Parse(reqs[0].URL)
	require.NoError(t, err)

	p := u.Path + "?" + u.RawQuery
	key := strings.TrimPrefix(u.Path, storageURLPath)

	badQuery := localstore.CipherQueryParam + "=" + base64.URLEncoding.EncodeToString([]byte("my_bad"))

	testCases := []struct {
		name         string
		path         string
		body         []byte
		expectedCode int
		key          string
	}{
		{name: "normal", path: p, key: key, body: []byte("test me"), expectedCode: http.StatusOK},
		{name: "bad_query", path: "/storage/testy?" + badQuery, body: []byte("test me"), expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, server.NewRequest(t, http.MethodPut, tc.path, bytes.NewReader(tc.body)))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &StoragePutOk{})
			if resp.Code == http.StatusOK {
				fileBytes, err := os.ReadFile(path.Join(routesConfig.StorageConfig.LocalStoreConfig.BaseBath, tc.key))
				require.NoError(err)
				require.Equal(string(tc.body), string(fileBytes))
			}
		})
	}
}
