package api

import (
	"bytes"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var notesServer txServer

func init() {
	notesServer = server.WithPathPrefix("/notes")
}

func TestRoutes_NoteCreate(t *testing.T) {
	testName := "TestRoutes_NoteCreate"

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{name: "basic", expectedCode: http.StatusOK},
		{name: "valid", expectedCode: http.StatusOK},
		{name: "invalid", expectedCode: http.StatusUnprocessableEntity},
	}
	for _, tc := range testCases {
		tc := tc
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
