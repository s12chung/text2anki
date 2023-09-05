package api

import (
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/util/chiutil"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
)

func idFromRequest(r *http.Request) (int64, *jhttp.HTTPError) {
	id, err := chiutil.ParamID(r, "id")
	if err != nil {
		return 0, jhttp.Error(http.StatusNotFound, err)
	}
	return id, nil
}

func txQsFromRequest(r *http.Request) (db.TxQs, *jhttp.HTTPError) {
	tx, err := reqtx.ContextTx(r)
	if err != nil {
		return db.TxQs{}, err
	}
	txQs, ok := tx.(db.TxQs) // type matches TxPool.GetTx
	if !ok {
		return db.TxQs{}, jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to db.TxQs fail"))
	}
	return txQs, nil
}

func extractAndValidate(r *http.Request, req any) *jhttp.HTTPError {
	if httpError := jhttp.ExtractJSON(r, req); httpError != nil {
		return httpError
	}
	result := firm.Validate(req)
	if !result.IsValid() {
		return jhttp.Error(http.StatusUnprocessableEntity, fmt.Errorf(result.ErrorMap().String()))
	}
	return nil
}

// TxPool is the default Pool for transactions
type TxPool struct{}

// GetTx returns a new transaction - returned type matches Routes.txQs()
func (t TxPool) GetTx(r *http.Request) (reqtx.Tx, error) {
	return db.NewTxQs(r.Context(), db.WriteOpts())
}
