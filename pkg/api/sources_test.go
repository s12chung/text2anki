package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
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
			testModelResponse(t, resp, testName, tc.name, &db.SourceSerialized{})
		})
	}
}

func TestRoutes_SourceUpdate(t *testing.T) {
	testName := "TestRoutes_SourceUpdate"

	testCases := []struct {
		name         string
		newName      string
		expectedCode int
	}{
		{name: "basic", newName: "new_name.txt", expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
	}

	created, err := db.Qs().SourceCreate(context.Background(), testdb.SourceSerializedsT(t)[1].ToSourceCreateParams())
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			reqBody := test.JSON(t, SourceUpdateRequest{Name: tc.newName})
			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodPatch, idPath("", created.ID), bytes.NewReader(reqBody)))
			require.Equal(tc.expectedCode, resp.Code)

			sourceSerialized := db.SourceSerialized{}
			fixtureFile := testModelResponse(t, resp, testName, tc.name, &sourceSerialized)

			if resp.Code == http.StatusOK {
				source, err := db.Qs().SourceGet(context.Background(), sourceSerialized.ID)
				require.NoError(err)
				fixture.CompareRead(t, fixtureFile, fixture.JSON(t, source.ToSourceSerialized().StaticCopy()))
			}
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
		{name: "empty", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			reqBody := test.JSON(t, sourcePostReqFromFile(t, testName, tc.name+".txt"))
			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodPost, "", bytes.NewReader(reqBody)))
			require.Equal(tc.expectedCode, resp.Code)

			sourceSerialized := db.SourceSerialized{}
			fixtureFile := testModelResponse(t, resp, testName, tc.name, &sourceSerialized)

			if resp.Code == http.StatusOK {
				source, err := db.Qs().SourceGet(context.Background(), sourceSerialized.ID)
				require.NoError(err)
				fixture.CompareRead(t, fixtureFile, fixture.JSON(t, source.ToSourceSerialized().StaticCopy()))
			}
		})
	}
}

func sourcePostReqFromFile(t *testing.T, testName, name string) *SourceCreateRequest {
	s := string(test.Read(t, fixture.JoinTestData(testName, name)))
	split := strings.Split(s, text.SplitDelimiter)
	if len(split) == 1 {
		return &SourceCreateRequest{Text: s}
	}
	return &SourceCreateRequest{Text: split[0], Translation: split[1]}
}

func TestRoutes_SourceDestroy(t *testing.T) {
	testName := "TestRoutes_SourceDestroy"

	created, err := db.Qs().SourceCreate(context.Background(), testdb.SourceSerializedsT(t)[1].ToSourceCreateParams())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "normal", path: idPath("", created.ID), expectedCode: http.StatusOK},
		{name: "invalid_id", path: "/9999", expectedCode: http.StatusNotFound},
		{name: "not_a_number", path: "/nan", expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodDelete, tc.path, nil))
			require.Equal(tc.expectedCode, resp.Code)

			testModelResponse(t, resp, testName, tc.name, &db.SourceSerialized{})
			if resp.Code == http.StatusOK {
				_, err := db.Qs().SourceGet(context.Background(), created.ID)
				require.Equal(fmt.Errorf("sql: no rows in result set"), err)
			}
		})
	}
}
