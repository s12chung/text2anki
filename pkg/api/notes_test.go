package api

import (
	"bytes"
	"net/http"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/archive/ziputil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var notesServer txServer

func init() {
	notesServer = server.WithPathPrefix("/notes")
}

func TestRoutes_NotesIndex(t *testing.T) {
	testName := "TestRoutes_NotesIndex"
	testIndex[db.Note](t, notesServer, testName, "notes")
}

func TestRoutes_NoteCreate(t *testing.T) {
	testName := "TestRoutes_NoteCreate"

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{name: "full", expectedCode: http.StatusOK},
		{name: "valid", expectedCode: http.StatusOK},
		{name: "invalid", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			txQs := testdb.TxQs(t, db.WriteOpts())

			reqBody := fixture.Read(t, path.Join(testName, tc.name+".json"))
			req := notesServer.NewTxRequest(t, txQs, http.MethodPost, "", bytes.NewReader(reqBody))
			resp := test.HTTPDo(t, req)
			resp.EqualCode(t, tc.expectedCode)

			note := db.Note{}
			fixtureFile := testModelResponse(t, resp, testName, tc.name, &note)
			if resp.Code != http.StatusOK {
				return
			}

			note, err := txQs.NoteGet(txQs.Ctx(), note.ID)
			require.NoError(err)
			fixture.CompareRead(t, fixtureFile, fixture.JSON(t, note.StaticCopy()))
		})
	}
}

func TestRoutes_NotesDownload(t *testing.T) {
	require := require.New(t)
	testName := "TestRoutes_NotesDownload"
	txQs := testdb.TxQs(t, db.WriteOpts())

	req := notesServer.NewTxRequest(t, txQs, http.MethodGet, "/download", nil)
	resp := test.HTTPDo(t, req)
	resp.EqualCode(t, http.StatusOK)
	result := resp.Result()
	require.Equal("attachment; filename=text2anki-"+testUUID+".zip", result.Header.Get("Content-Disposition"))
	require.NoError(result.Body.Close())

	paths := []string{
		"files/",
		"files/t2a-꽃길만 걷게 해줄게요.mp3",
		"files/t2a-모자람 없이 주신 사랑이 과분하다 느낄 때쯤 난 어른이 됐죠.mp3",
		"files/t2a-여길 봐 예쁘게 피었으니까.mp3",
		"text2anki.csv",
	}
	require.NoError(ziputil.CompareContents(resp.Body.Bytes(), paths, func(path string, contents []byte) {
		fixture.CompareReadOrUpdate(t, filepath.Join(testName, path), contents) //nolint:gosec // for testing
	}))

	notes, err := txQs.NotesDownloaded(txQs.Ctx())
	require.NoError(err)
	require.Nil(notes)
}
