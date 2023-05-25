package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestRoutes_SourceList(t *testing.T) {
	require := require.New(t)
	testName := "TestRoutes_SourceList"

	resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, "", nil))
	require.Equal(http.StatusOK, resp.Code)
	fixture.CompareReadOrUpdate(t, testName+".json", test.StaticCopySlice(t, resp.Body.Bytes(), &[]db.SourceSerialized{}))
}

func TestRoutes_SourceGet(t *testing.T) {
	testName := "TestRoutes_SourceGet"
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

			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, tc.path, nil))
			require.Equal(tc.expectedCode, resp.Code)

			jsonBody := test.StaticCopyOrIndent(t, resp.Code, resp.Body.Bytes(), &db.SourceSerialized{})
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), jsonBody)
		})
	}
}

func TestRoutes_SourceCreate(t *testing.T) {
	testName := "TestRoutes_SourceCreate"
	test.CISkip(t, "can't run C environment in CI")

	require.NoError(t, DefaultRoutes.Setup())
	defer func() {
		require.NoError(t, DefaultRoutes.Cleanup())
	}()

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
			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodPost, "", bytes.NewReader(reqBody)))
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

func TestRoutes_SourceDestroy(t *testing.T) {
	testName := "TestRoutes_SourceDestroy"

	created, err := db.Qs().SourceSerializedCreate(context.Background(), testdb.SourceSerializedsT(t)[1].TokenizedTexts)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "normal", path: fmt.Sprintf("/%v", created.ID), expectedCode: http.StatusOK},
		{name: "invalid_id", path: "/9999", expectedCode: http.StatusNotFound},
		{name: "not_a_number", path: "/nan", expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodDelete, tc.path, nil))
			require.Equal(tc.expectedCode, resp.Code)

			jsonBody := test.StaticCopyOrIndent(t, resp.Code, resp.Body.Bytes(), &db.SourceSerialized{})
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), jsonBody)

			if resp.Code == http.StatusOK {
				_, err := db.Qs().SourceGet(context.Background(), created.ID)
				require.Equal(fmt.Errorf("sql: no rows in result set"), err)
			}
		})
	}
}
