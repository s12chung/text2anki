package api

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/test"
)

const txtExt = ".txt"

func TestRoutes_PartCreateMulti(t *testing.T) {
	testName := "TestRoutes_PartCreateMulti"
	test.CISkip(t, "can't run C environment in CI")
	t.Parallel()

	mediaID := "b1234567-3456-9abc-d123-456789abcdef"
	mediaWithInfoID := "b47ac10b-58cc-4372-a567-0e02b2c3d479"
	setupSourceCreateMedia(t, mediaID)
	setupSourceCreateMediaWithInfo(t, mediaWithInfoID)

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
		{name: "multi_with_empty", partCount: 3, expectedCode: http.StatusOK},
		{name: "media", partCount: 3, prePartListID: mediaID, expectedCode: http.StatusOK},
		{name: "error", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty", expectedCode: http.StatusUnprocessableEntity},
		{name: "empty_parts", expectedCode: http.StatusUnprocessableEntity},
		{name: "51_parts", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			txQs := testdb.TxQs(t, db.WriteOpts())
			created := createdSource(t, txQs)

			body := PartCreateMultiRequest{PrePartListID: tc.prePartListID, Parts: sourceCreateRequestParts(t, tc.name, testName, tc.partCount)}
			req := sourcesServer.NewTxRequest(t, txQs, http.MethodPost, joinPath(created.ID, "parts", "multi"), bytes.NewReader(test.JSON(t, body)))
			testSourceResponse(t, req, txQs, testName, tc.name, tc.expectedCode)
		})
	}
}

func TestRoutes_PartCreate(t *testing.T) {
	testName := "TestRoutes_PartCreate"
	test.CISkip(t, "can't run C environment in CI")
	t.Parallel()

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{name: "basic", expectedCode: http.StatusOK},
		{name: "no_translation", expectedCode: http.StatusOK},
		{name: "empty", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			txQs := testdb.TxQs(t, db.WriteOpts())
			created := createdSource(t, txQs)

			body := sourceCreateRequestPartFromFile(t, testName, tc.name+txtExt)
			req := sourcesServer.NewTxRequest(t, txQs, http.MethodPost, joinPath(created.ID, "parts"), bytes.NewReader(test.JSON(t, body)))
			testSourceResponse(t, req, txQs, testName, tc.name, tc.expectedCode)
		})
	}
}

func TestRoutes_PartUpdate(t *testing.T) {
	testName := "TestRoutes_PartUpdate"
	test.CISkip(t, "can't run C environment in CI")
	t.Parallel()

	testCases := []struct {
		name         string
		index        int
		expectedCode int
	}{
		{name: "basic", expectedCode: http.StatusOK},
		{name: "no_translation", expectedCode: http.StatusOK},
		{name: "index_out", index: 1, expectedCode: http.StatusUnprocessableEntity},
		{name: "index_negative", index: -1, expectedCode: http.StatusUnprocessableEntity},
		{name: "empty", index: 1, expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			txQs := testdb.TxQs(t, db.WriteOpts())
			created := createdSource(t, txQs)

			body := sourceCreateRequestPartFromFile(t, testName, tc.name+txtExt)
			req := sourcesServer.NewTxRequest(t, txQs, http.MethodPatch, joinPath(created.ID, "parts", tc.index), bytes.NewReader(test.JSON(t, body)))
			testSourceResponse(t, req, txQs, testName, tc.name, tc.expectedCode)
		})
	}
}

func TestRoutes_PartDestroy(t *testing.T) {
	testName := "TestRoutes_PartDestroy"
	test.CISkip(t, "can't run C environment in CI")
	t.Parallel()

	testCases := []struct {
		name         string
		index        int
		expectedCode int
	}{
		{name: "basic", expectedCode: http.StatusOK},
		{name: "middle", index: 1, expectedCode: http.StatusOK},
		{name: "end", index: 2, expectedCode: http.StatusOK},
		{name: "index_out", index: 3, expectedCode: http.StatusUnprocessableEntity},
		{name: "index_negative", index: -1, expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			t.Parallel()

			txQs := testdb.TxQs(t, db.WriteOpts())
			model := testdb.SourceStructureds().ModelsT(t)[1]
			tokenizedTexts, err := routes.TextTokenizer.TokenizedTexts(txQs.Ctx(), "내가 말한 '혼자'는 너랑 같이 있어야 왼성되거든.", "")
			require.NoError(err)
			model.Parts = append(model.Parts, append(testdb.SourceStructureds().ModelsT(t)[0].Parts, db.SourcePart{TokenizedTexts: tokenizedTexts})...)
			created, err := txQs.SourceCreate(txQs.Ctx(), model.CreateParams())
			require.NoError(err)

			req := sourcesServer.NewTxRequest(t, txQs, http.MethodDelete, joinPath(created.ID, "parts", tc.index), nil)
			testSourceResponse(t, req, txQs, testName, tc.name, tc.expectedCode)
		})
	}
}
