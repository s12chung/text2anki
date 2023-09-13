package api

import (
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

func init() {
	httptyped.RegisterType(dictionary.Term{})
}

// TermsSearch searches for the terms
func (rs Routes) TermsSearch(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	query := r.URL.Query().Get("query")
	posQuery := r.URL.Query().Get("pos")
	pos, err := lang.ToPartOfSpeech(posQuery)
	if err != nil {
		return nil, jhttp.Error(http.StatusUnprocessableEntity, err)
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
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
