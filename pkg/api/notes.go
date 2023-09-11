package api

import (
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

func init() {
	httptyped.RegisterType(db.Note{})
}

// NotesIndex shows lists all the notes
func (rs Routes) NotesIndex(_ *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	return jhttp.ReturnModelOr500(func() (any, error) {
		return txQs.NotesIndex(txQs.Ctx())
	})
}

// NoteCreate creates a new note
func (rs Routes) NoteCreate(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	req := db.NoteCreateParams{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		return txQs.NoteCreate(txQs.Ctx(), req)
	})
}
