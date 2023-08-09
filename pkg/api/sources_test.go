package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/models"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var sourcesServer test.Server

func init() {
	sourcesServer = server.WithPathPrefix("/sources")
}

func TestRoutes_SourceIndex(t *testing.T) {
	testName := "TestRoutes_SourceIndex"
	resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, "", nil))
	testModelsResponse(t, resp, testName, "", &[]db.SourceSerialized{})
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
		tc := tc
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

	created, err := db.Qs().SourceCreate(context.Background(), models.SourceSerializedsMust()[1].CreateParams())
	require.NoError(t, err)

	for _, tc := range testCases {
		tc := tc
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

	prePartListID := testUUID
	setupSourceCreateMedia(t, prePartListID)
	require.NoError(t, routes.Setup())
	defer func() {
		require.NoError(t, routes.Cleanup())
	}()

	testCases := []struct {
		name          string
		partCount     int
		prePartListID string
		expectedCode  int
	}{
		{name: "split", expectedCode: http.StatusOK},
		{name: "no_translation", expectedCode: http.StatusOK},
		{name: "weave", expectedCode: http.StatusOK},
		{name: "multi_part", partCount: 2, expectedCode: http.StatusOK},
		{name: "media", partCount: 3, prePartListID: prePartListID, expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty_parts", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			body := SourceCreateRequest{PrePartListID: tc.prePartListID}
			if tc.partCount == 0 {
				body.Parts = []SourceCreateRequestPart{sourcePartFromFile(t, testName, tc.name+".txt")}
			} else {
				body.Parts = make([]SourceCreateRequestPart, tc.partCount)
				for i := 0; i < tc.partCount; i++ {
					body.Parts[i] = sourcePartFromFile(t, testName, tc.name+strconv.Itoa(i)+".txt")
				}
			}
			if tc.name == "empty_parts" {
				body.Parts = []SourceCreateRequestPart{}
			}

			reqBody := test.JSON(t, body)
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

func setupSourceCreateMedia(t *testing.T, prePartListID string) {
	baseKey := storage.BaseKey(sourcesTable, partsColumn, prePartListID)

	for i := 0; i < 3; i++ {
		err := routes.Storage.Storer.Store(baseKey+".PreParts["+strconv.Itoa(i)+"].Image.txt", bytes.NewReader([]byte("image"+strconv.Itoa(i))))
		require.NoError(t, err)
	}
	err := routes.Storage.Storer.Store(baseKey+".PreParts[0].Audio.txt", bytes.NewReader([]byte("audio0!")))
	require.NoError(t, err)
}

func sourcePartFromFile(t *testing.T, testName, name string) SourceCreateRequestPart {
	s := string(test.Read(t, fixture.JoinTestData(testName, name)))
	split := strings.Split(s, "===")
	part := SourceCreateRequestPart{Text: s}
	if len(split) == 2 {
		part = SourceCreateRequestPart{Text: split[0], Translation: split[1]}
	}
	return part
}

func TestRoutes_SourceDestroy(t *testing.T) {
	testName := "TestRoutes_SourceDestroy"

	created, err := db.Qs().SourceCreate(context.Background(), models.SourceSerializedsMust()[1].CreateParams())
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
		tc := tc
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
