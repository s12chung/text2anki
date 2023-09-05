package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/util/chiutil"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

func idFromRequest(r *http.Request) (int64, *jhttp.HTTPError) {
	id, err := chiutil.ParamID(r, "id")
	if err != nil {
		return 0, jhttp.Error(http.StatusNotFound, err)
	}
	return id, nil
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

const (
	txReadOnly config.TxMode = iota
	txWritable
)

var txModeToOpts = map[config.TxMode]sql.TxOptions{
	txReadOnly: {ReadOnly: true},
	txWritable: {},
}

// GetTx returns a new transaction
func (t TxPool) GetTx(r *http.Request, mode config.TxMode) (db.TxQs, error) {
	opts, exists := txModeToOpts[mode]
	if !exists {
		return db.TxQs{}, fmt.Errorf("config.TxMode does not exist: %v", mode)
	}
	return db.NewTxQs(r.Context(), &opts)
}
