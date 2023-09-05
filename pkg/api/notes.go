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

// NoteCreate creates a new note
func (rs Routes) NoteCreate(r *http.Request) (any, *jhttp.HTTPError) {
	req := db.NoteCreateParams{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		return txQs.NoteCreate(txQs.Ctx(), req)
	})
}
