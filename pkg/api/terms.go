package api

import (
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/httptyped"
)

func init() {
	httptyped.RegisterType(dictionary.Term{})
}

var posTypes = lang.PartOfSpeechTypes()

// TermsSearch searches for the terms
func (rs Routes) TermsSearch(r *http.Request) (any, *httputil.HTTPError) {
	query := r.URL.Query().Get("query")
	posQuery := r.URL.Query().Get("pos")
	pos, found := posTypes[posQuery]
	if !found {
		return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("pos is invalid: '%v'", posQuery))
	}

	txQs, httpErr := rs.txQs(r)
	if httpErr != nil {
		return nil, httpErr
	}
	return httputil.ReturnModelOr500(func() (any, error) {
		termsSearchRow, err := txQs.TermsSearch(txQs.Ctx(), query, pos)
		if err != nil {
			return nil, err
		}
		terms := make([]dictionary.Term, len(termsSearchRow))
		for i, row := range termsSearchRow {
			term, err := row.Term.DictionaryTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		}
		return terms, err
	})
}
