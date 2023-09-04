package api

import (
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/util/httputil"
	"github.com/s12chung/text2anki/pkg/util/httputil/reqtx"
)

func extractAndValidate(r *http.Request, req any) *httputil.HTTPError {
	if httpError := httputil.ExtractJSON(r, req); httpError != nil {
		return httpError
	}
	result := firm.Validate(req)
	if !result.IsValid() {
		return httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf(result.ErrorMap().String()))
	}
	return nil
}

func (rs Routes) txQs(r *http.Request) (db.TxQs, *httputil.HTTPError) {
	tx, err := rs.TxIntegrator.ContextTx(r)
	if err != nil {
		return db.TxQs{}, err
	}
	txQs, ok := tx.(db.TxQs) // type matches TxPool.GetTx
	if !ok {
		return db.TxQs{}, httputil.Error(http.StatusInternalServerError, fmt.Errorf("cast to db.TxQs fail"))
	}
	return txQs, nil
}

// TxPool is the default Pool for transactions
type TxPool struct{}

// GetTx returns a new transaction - returned type matches Routes.txQs()
func (t TxPool) GetTx(r *http.Request) (reqtx.Tx, error) {
	return db.NewTxQs(r.Context(), db.WriteOpts())
}
