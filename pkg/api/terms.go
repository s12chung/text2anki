package api

import (
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/httputil"
)

var posTypes = lang.PartOfSpeechTypes()

// TermsSearch searches for the terms
func (rs Routes) TermsSearch(r *http.Request) (any, int, error) {
	query := r.URL.Query().Get("query")
	posQuery := r.URL.Query().Get("pos")
	pos, found := posTypes[posQuery]
	if !found {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf("pos is invalid: '%v'", posQuery)
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		termsSearchRow, err := db.Qs().TermsSearch(r.Context(), query, pos)
		terms := make([]db.Term, len(termsSearchRow))
		for i, row := range termsSearchRow {
			terms[i] = row.Term
		}
		return terms, err
	})
}
