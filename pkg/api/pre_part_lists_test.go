package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/test"
)

var prePartListServer test.Server

func init() {
	prePartListServer = server.WithPathPrefix("/sources/pre_part_lists")
}

type signPrePartListResponse struct {
	PrePartListSignResponse
}

func replaceCipherQueryParam(urlString string) string {
	if urlString == "" {
		return ""
	}

	u, err := url.Parse(urlString)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	u.RawQuery = url.Values{localstore.CipherQueryParam: []string{"testy"}}.Encode()
	return u.String()
}

func (s signPrePartListResponse) StaticCopy() any {
	a := s
	a.ID = testUUID
	for i, prePart := range s.PreParts {
		if prePart.ImageRequest != nil {
			prePart.ImageRequest.URL = replaceCipherQueryParam(prePart.ImageRequest.URL)
		}
		if prePart.AudioRequest != nil {
			prePart.AudioRequest.URL = replaceCipherQueryParam(prePart.AudioRequest.URL)
		}
		a.PreParts[i] = prePart
	}
	return a
}

func TestRoutes_PrePartListSign(t *testing.T) {
	testName := "TestRoutes_PrePartListSign"

	testCases := []struct {
		name         string
		req          PrePartListSignRequest
		expectedCode int
	}{
		{name: "one", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg"}}},
			expectedCode: http.StatusOK},
		{name: "many", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg"}, {ImageExt: ".png"}, {ImageExt: ".jpeg"}}},
			expectedCode: http.StatusOK},
		{name: "mixed", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg", AudioExt: ".mp3"}, {AudioExt: ".mp3"}, {ImageExt: ".jpeg"}}},
			expectedCode: http.StatusOK},
		{name: "none", req: PrePartListSignRequest{}, expectedCode: http.StatusUnprocessableEntity},
		{name: "invalid", req: PrePartListSignRequest{
			PreParts: []PrePartSignRequest{{ImageExt: ".jpg", AudioExt: ".mp3"}, {ImageExt: ".waka"}}},
			expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, prePartListServer.NewRequest(t, http.MethodPost, "/sign", bytes.NewReader(test.JSON(t, tc.req))))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &signPrePartListResponse{})
		})
	}
}

func TestRoutes_PrePartListGet(t *testing.T) {
	testName := "TestRoutes_PrePartListGet"

	id := "my_id"
	baseKey := storage.BaseKey(sourcesTable, partsColumn, id)
	for i := 0; i < 2; i++ {
		err := routes.Storage.Storer.Store(baseKey+".PreParts["+strconv.Itoa(i)+"].Image.txt", bytes.NewReader([]byte("image"+strconv.Itoa(i))))
		require.NoError(t, err)
	}
	err := routes.Storage.Storer.Store(baseKey+".PreParts[0].Audio.txt", bytes.NewReader([]byte("audio0!")))
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

			resp := test.HTTPDo(t, prePartListServer.NewRequest(t, http.MethodGet, "/"+tc.id, nil))
			require.Equal(tc.expectedCode, resp.Code)
			testModelResponse(t, resp, testName, tc.name, &PrePartList{})
		})
	}
}
