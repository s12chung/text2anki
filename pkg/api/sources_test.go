package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestRoutes_SourceGet(t *testing.T) {
	testName := "TestRoutes_SourceGet"

	r := chi.NewRouter()
	r.Route("/{sourceID}", func(r chi.Router) {
		r.Use(httputil.RequestWrap(SourceCtx))
		r.Get("/", httputil.RespondJSONWrap(DefaultRoutes.SourceGet))
	})

	server := httptest.NewServer(r)
	defer server.Close()

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "normal", path: "/1", expectedCode: http.StatusOK},
		{name: "invalid_id", path: "/9999", expectedCode: http.StatusNotFound},
		{name: "not_a_number", path: "/nan", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			req, err := http.NewRequest(http.MethodGet, server.URL+tc.path, nil)
			require.NoError(err)

			resp := test.HTTPDo(t, req)
			require.Equal(tc.expectedCode, resp.Code)

			jsonBody := test.StaticCopyOrIndent(t, resp.Code, resp.Body.Bytes(), &db.SourceSerialized{})
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), jsonBody)
		})
	}
}

func TestRoutes_SourcePost(t *testing.T) {
	test.CISkip(t, "can't run C environment in CI")
	require.NoError(t, DefaultRoutes.Setup())
	defer func() {
		require.NoError(t, DefaultRoutes.Cleanup())
	}()

	testName := "TestRoutes_SourcePost"

	handlerFunc := httputil.RespondJSONWrap(DefaultRoutes.SourcePost)

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{name: "split", expectedCode: http.StatusOK},
		{name: "weave", expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			reqBody := test.JSON(t, sourcePostReqFromFile(t, testName, tc.name+".txt"))
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
			resp := httptest.NewRecorder()

			handlerFunc(resp, req)
			require.Equal(tc.expectedCode, resp.Code)

			sourceSerialized := db.SourceSerialized{}
			jsonBody := test.StaticCopyOrIndent(t, resp.Code, resp.Body.Bytes(), &sourceSerialized)
			fixtureFile := path.Join(testName, tc.name+"_response.json")
			fixture.CompareReadOrUpdate(t, fixtureFile, jsonBody)

			if resp.Code == http.StatusOK {
				source, err := db.Qs().SourceGet(context.Background(), sourceSerialized.ID)
				require.NoError(err)
				fixture.CompareRead(t, fixtureFile, fixture.JSON(t, source.ToSourceSerialized().StaticCopy()))
			}
		})
	}
}

func sourcePostReqFromFile(t *testing.T, testName, name string) *SourcePostRequest {
	s := string(test.Read(t, fixture.JoinTestData(testName, name)))
	split := strings.Split(s, text.SplitDelimiter)
	if len(split) == 1 {
		return &SourcePostRequest{Text: s}
	}
	return &SourcePostRequest{Text: split[0], Translation: split[1]}
}
