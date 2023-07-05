package api

import (
	"bytes"
	"context"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var notesServer test.Server

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

			reqBody := fixture.Read(t, path.Join(testName, tc.name+".json"))
			resp := test.HTTPDo(t, notesServer.NewRequest(t, http.MethodPost, "", bytes.NewReader(reqBody)))
			require.Equal(tc.expectedCode, resp.Code)

			note := db.Note{}
			fixtureFile := testModelResponse(t, resp, testName, tc.name, &note)

			if resp.Code == http.StatusOK {
				note, err := db.Qs().NoteGet(context.Background(), note.ID)
				require.NoError(err)
				fixture.CompareRead(t, fixtureFile, fixture.JSON(t, note.StaticCopy()))
			}
		})
	}
}