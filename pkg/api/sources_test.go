package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var sourcesServer txServer

func init() {
	sourcesServer = server.WithPathPrefix("/sources")
}
func createdSource(t *testing.T, txQs db.TxQs) db.Source {
	require := require.New(t)
	created, err := txQs.SourceCreate(txQs.Ctx(), testdb.SourceStructureds().ModelsT(t)[1].CreateParams())
	require.NoError(err)
	return created
}

func TestRoutes_SourcesIndex(t *testing.T) {
	testName := "TestRoutes_SourcesIndex"
	resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, "", nil))
	testModelsResponse[db.SourceStructured](t, resp, testName, "", nil)
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
			resp := test.HTTPDo(t, sourcesServer.NewRequest(t, http.MethodGet, tc.path, nil))
			resp.EqualCode(t, tc.expectedCode)
			testModelResponse(t, resp, testName, tc.name, &db.SourceStructured{})
		})
	}
}

func TestRoutes_SourceCreate(t *testing.T) {
	testName := "TestRoutes_SourceCreate"
	test.CISkip(t, "can't run C environment in CI")
	t.Parallel()

	mediaID := "a1234567-3456-9abc-d123-456789abcdef"
	mediaWithInfoID := "a47ac10b-58cc-4372-a567-0e02b2c3d479"
	setupSourceCreateMedia(t, mediaID)
	setupSourceCreateMediaWithInfo(t, mediaWithInfoID)

	testCases := []struct {
		name          string
		sName         string
		reference     string
		partCount     int
		prePartListID string
		expectedCode  int
	}{
		{name: "split", expectedCode: http.StatusOK},
		{name: "split_with_reference", reference: "split_with_reference.txt", expectedCode: http.StatusOK},
		{name: "split_with_name_and_ref", sName: "some_name", reference: "split_with_name_and_ref.txt", expectedCode: http.StatusOK},
		{name: "no_translation", expectedCode: http.StatusOK},
		{name: "weave", expectedCode: http.StatusOK},
		{name: "multi_part", partCount: 2, expectedCode: http.StatusOK},
		{name: "multi_with_empty", partCount: 3, expectedCode: http.StatusOK},
		{name: "media", partCount: 3, prePartListID: mediaID, expectedCode: http.StatusOK},
		{name: "media_with_info", partCount: 2, prePartListID: mediaWithInfoID, expectedCode: http.StatusOK},
		{name: "media_with_info_name_ref", partCount: 2, prePartListID: mediaWithInfoID,
			sName: "override name", reference: "override.txt", expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty_parts", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			txQs := testdb.TxQs(t, db.WriteOpts())

			body := SourceCreateRequest{
				Name:      tc.sName,
				Reference: tc.reference,
				PartCreateMultiRequest: PartCreateMultiRequest{
					PrePartListID: tc.prePartListID,
					Parts:         sourceCreateRequestParts(t, tc.name, testName, tc.partCount),
				},
			}
			req := sourcesServer.NewTxRequest(t, txQs, http.MethodPost, "", bytes.NewReader(test.JSON(t, body)))
			testSourceResponse(t, req, txQs, testName, tc.name, tc.expectedCode)
		})
	}
}

func TestRoutes_SourceUpdate(t *testing.T) {
	testName := "TestRoutes_SourceUpdate"
	t.Parallel()

	testCases := []struct {
		name         string
		newName      string
		reference    string
		expectedCode int
	}{
		{name: "basic", newName: "new_name", expectedCode: http.StatusOK},
		{name: "with_reference", newName: "new_name", reference: "new_ref.txt", expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			txQs := testdb.TxQs(t, db.WriteOpts())
			created := createdSource(t, txQs)

			body := SourceUpdateRequest{Name: tc.newName, Reference: tc.reference}
			req := sourcesServer.NewTxRequest(t, txQs, http.MethodPatch, joinPath(created.ID), bytes.NewReader(test.JSON(t, body)))
			testSourceResponse(t, req, txQs, testName, tc.name, tc.expectedCode)
		})
	}
}

func TestRoutes_SourceDestroy(t *testing.T) {
	testName := "TestRoutes_SourceDestroy"
	t.Parallel()

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "normal", expectedCode: http.StatusOK},
		{name: "invalid_id", path: "/9999", expectedCode: http.StatusNotFound},
		{name: "not_a_number", path: "/nan", expectedCode: http.StatusNotFound},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			t.Parallel()

			txQs := testdb.TxQs(t, db.WriteOpts())
			created := createdSource(t, txQs)
			if tc.path == "" {
				tc.path = joinPath(created.ID)
			}

			resp := test.HTTPDo(t, sourcesServer.NewTxRequest(t, txQs, http.MethodDelete, tc.path, nil))
			resp.EqualCode(t, tc.expectedCode)

			testModelResponse(t, resp, testName, tc.name, &db.SourceStructured{})
			if resp.Code != http.StatusOK {
				return
			}

			_, err := txQs.SourceGet(txQs.Ctx(), created.ID)
			require.Equal(fmt.Errorf("sql: no rows in result set"), err)
		})
	}
}

func setupSourceCreateMediaWithInfo(t *testing.T, prePartListID string) {
	setupSourceCreateMedia(t, prePartListID)

	info := db.PrePartInfo{Name: "test name", Reference: "https://www.testref.com"}
	baseKey := storage.BaseKey(db.SourcesTable, db.PartsColumn, prePartListID)
	err := routes.Storage.Storer.Store(baseKey+".Info.json", bytes.NewReader(test.JSON(t, info)))
	require.NoError(t, err)
}

func setupSourceCreateMedia(t *testing.T, prePartListID string) {
	baseKey := storage.BaseKey(db.SourcesTable, db.PartsColumn, prePartListID)

	for i := 0; i < 3; i++ {
		err := routes.Storage.Storer.Store(baseKey+".PreParts["+strconv.Itoa(i)+"].Image.txt", bytes.NewReader([]byte("image"+strconv.Itoa(i))))
		require.NoError(t, err)
	}
	err := routes.Storage.Storer.Store(baseKey+".PreParts[0].Audio.txt", bytes.NewReader([]byte("audio0!")))
	require.NoError(t, err)
}

func sourceCreateRequestParts(t *testing.T, caseName, testName string, partCount int) []PartCreateMultiRequestPart {
	if partCount == 0 {
		partCount = 1
	}
	parts := make([]PartCreateMultiRequestPart, partCount)
	switch caseName {
	case "empty":
		parts[0].Text = "  "
	case "empty_parts":
		parts = []PartCreateMultiRequestPart{}
	default:
		for i := 0; i < partCount; i++ {
			parts[i] = sourceCreateRequestPartFromFile(t, testName, caseName+strconv.Itoa(i)+".txt")
		}
	}
	return parts
}

func sourceCreateRequestPartFromFile(t *testing.T, testName, name string) PartCreateMultiRequestPart {
	s := string(test.Read(t, fixture.JoinTestData(testName, name)))
	split := strings.Split(s, "===")
	part := PartCreateMultiRequestPart{Text: s}
	if len(split) == 2 {
		part = PartCreateMultiRequestPart{Text: split[0], Translation: split[1]}
	}
	return part
}

func testSourceResponse(t *testing.T, req *http.Request, txQs db.TxQs, testName, name string, expectedCode int) {
	require := require.New(t)

	resp := test.HTTPDo(t, req)
	resp.EqualCode(t, expectedCode)

	sourceStructured := db.SourceStructured{}
	fixtureFile := testModelResponse(t, resp, testName, name, &sourceStructured)
	if resp.Code != http.StatusOK {
		return
	}

	source, err := txQs.SourceGet(txQs.Ctx(), sourceStructured.ID)
	require.NoError(err)
	sourceStructured = source.ToSourceStructured()
	sourceStructured.PrepareSerialize()
	fixture.CompareRead(t, fixtureFile, fixture.JSON(t, sourceStructured.StaticCopy()))
}
