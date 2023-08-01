package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/test"
)

var prePartsServer test.Server

func init() {
	prePartsServer = server.WithPathPrefix("/sources/pre_parts")
}

type signPrePartsResponse struct {
	PrePartsSignResponse
}

var testUUID = "123e4567-e89b-12d3-a456-426614174000"
var uuidRegexp = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)

func (s signPrePartsResponse) StaticCopy() any {
	a := s
	a.ID = testUUID
	for i, req := range s.Requests {
		u, err := url.Parse(req.URL)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		u.RawQuery = url.Values{localstore.CipherQueryParam: []string{"testy"}}.Encode()
		req.URL = uuidRegexp.ReplaceAllString(u.String(), testUUID)
		a.Requests[i] = req
	}
	return a
}

func TestRoutes_PrePartsSign(t *testing.T) {
	testName := "TestRoutes_PrePartsSign"

	testCases := []struct {
		name         string
		queryParams  string
		expectedCode int
	}{
		{name: "one", queryParams: "exts=.png", expectedCode: http.StatusOK},
		{name: "many", queryParams: "exts=.jpg&exts=.png&exts=.jpeg", expectedCode: http.StatusOK},
		{name: "array", queryParams: "exts[0]=.jpg&exts[1]=.png&exts[2]=.jpeg", expectedCode: http.StatusUnprocessableEntity},
		{name: "comma", queryParams: "exts=.jpeg,.png", expectedCode: http.StatusUnprocessableEntity},
		{name: "none", queryParams: "", expectedCode: http.StatusUnprocessableEntity},
		{name: "invalid", queryParams: "exts=.jpg&exts=.png&exts=.waka", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, prePartsServer.NewRequest(t, http.MethodGet, "/sign?"+tc.queryParams, nil))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &signPrePartsResponse{})
		})
	}
}

func TestRoutes_PrePartsGet(t *testing.T) {
	testName := "TestRoutes_PrePartsGet"

	id := "my_id"
	idPath := path.Join(sourcesTable, partsColumn, id)
	err := routes.Storage.Storer.Store(path.Join(idPath, "blah.txt"), bytes.NewReader([]byte("my_blah")))
	require.NoError(t, err)
	err = routes.Storage.Storer.Store(path.Join(idPath, "again_me.txt"), bytes.NewReader([]byte("again!")))
	require.NoError(t, err)

	testCases := []struct {
		name         string
		id           string
		expectedCode int
	}{
		{name: "many", id: id, expectedCode: http.StatusOK},
		{name: "none", id: "does_not_exist", expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, prePartsServer.NewRequest(t, http.MethodGet, "/"+tc.id, nil))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &PreParts{})
		})
	}
}
