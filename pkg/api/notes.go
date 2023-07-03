package api

import (
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/httputil"
)

// NoteCreate creates a new note
func (rs Routes) NoteCreate(r *http.Request) (any, int, error) {
	req := db.NoteCreateParams{}
	if code, err := extractAndValidate(r, &req); err != nil {
		return nil, code, err
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return db.Qs().NoteCreate(r.Context(), req)
	})
}
