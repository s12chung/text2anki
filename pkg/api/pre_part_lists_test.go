package api

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/extractor/extractortest"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/logg"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var prePartListServer txServer

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
		slog.Error("api.replaceCipherQueryParam()", logg.Err(err))
		os.Exit(-1)
	}
	u.RawQuery = url.Values{localstore.CipherQueryParam: []string{"testy"}}.Encode()
	return u.String()
}

func (s signPrePartListResponse) StaticCopy() signPrePartListResponse {
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
			resp := test.HTTPDo(t, prePartListServer.NewRequest(t, http.MethodPost, "/sign", bytes.NewReader(test.JSON(t, tc.req))))
			resp.EqualCode(t, tc.expectedCode)
			testModelResponse(t, resp, testName, tc.name, &signPrePartListResponse{})
		})
	}
}

func TestRoutes_PrePartListGet(t *testing.T) {
	testName := "TestRoutes_PrePartListGet"

	id := "my_id"
	baseKey := storage.BaseKey(db.SourcesTable, db.PartsColumn, id)
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
			resp := test.HTTPDo(t, prePartListServer.NewRequest(t, http.MethodGet, "/"+tc.id, nil))
			resp.EqualCode(t, tc.expectedCode)
			testModelResponse(t, resp, testName, tc.name, &db.PrePartListURL{})
		})
	}
}

func TestRoutes_PrePartListVerify(t *testing.T) {
	testName := "TestRoutes_PrePartListVerify"
	testCases := []struct {
		name         string
		text         string
		expectedType string
	}{
		{name: "matched", text: extractortest.VerifyString, expectedType: extractorType},
		{name: "not_matched", text: "does not match", expectedType: ""},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := PrePartListVerifyRequest{Text: tc.text}
			resp := test.HTTPDo(t, prePartListServer.NewRequest(t, http.MethodPost, "/verify", bytes.NewReader(test.JSON(t, req))))
			resp.EqualCode(t, http.StatusOK)
			testModelResponse(t, resp, testName, tc.name, &PrePartListVerifyResponse{})
		})
	}
}

func TestRoutes_PrePartListCreate(t *testing.T) {
	testName := "TestRoutes_PrePartListCreate"
	testCases := []struct {
		name         string
		req          PrePartListCreateRequest
		expectedCode int
	}{
		{name: "matched", req: PrePartListCreateRequest{ExtractorType: extractorType, Text: extractortest.VerifyString}, expectedCode: http.StatusOK},
		{name: "bad_string", req: PrePartListCreateRequest{ExtractorType: extractorType, Text: "bad_string"},
			expectedCode: http.StatusUnprocessableEntity},
		{name: "bad_type", req: PrePartListCreateRequest{ExtractorType: "bad_type", Text: extractortest.VerifyString},
			expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			resp := test.HTTPDo(t, prePartListServer.NewRequest(t, http.MethodPost, "/", bytes.NewReader(test.JSON(t, tc.req))))
			resp.EqualCode(t, tc.expectedCode)
			prePartListResp := PrePartListCreateResponse{}
			testModelResponse(t, resp, testName, tc.name, &prePartListResp)

			if resp.Code != http.StatusOK {
				return
			}

			keyTree := db.PrePartList{}
			err := routes.Storage.DBStorage.KeyTree(db.SourcesTable, db.PartsColumn, prePartListResp.ID, &keyTree)
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+"_KeyTree.json"), fixture.JSON(t, keyTree))

			for _, prePart := range keyTree.PreParts {
				file, err := routes.Storage.DBStorage.Get(prePart.ImageKey)
				require.NoError(err)
				fileBytes, err := io.ReadAll(file)
				require.NoError(err)
				require.Equal("image_content", string(fileBytes))
			}
		})
	}
}
