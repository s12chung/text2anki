package api

import (
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(db.Note{})
}

// NoteCreate creates a new note
func (rs Routes) NoteCreate(r *http.Request) (any, *httputil.HTTPError) {
	req := db.NoteCreateParams{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		return txQs.NoteCreate(txQs.Ctx(), req)
	})
}
