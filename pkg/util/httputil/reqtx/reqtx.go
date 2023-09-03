// Package reqtx manages a database transaction per request
package reqtx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/pkg/util/httputil"
)

const txContextKey httputil.ContextKey = "reqtx.TxContext"

// Tx represents a transaction
type Tx interface {
	Finalize() error
	FinalizeError() error
}

// Pool gets tranactions to handle
type Pool interface {
	GetTx(r *http.Request) (Tx, error)
}

// Integrator integrates Pool to the request via the Router
type Integrator struct{ pool Pool }

// NewIntegrator returns a new Integrator
func NewIntegrator(pool Pool) Integrator { return Integrator{pool: pool} }

// SetTxContext sets the transaction on the request context
func (i Integrator) SetTxContext(r *http.Request) (*http.Request, *httputil.HTTPError) {
	tx, err := i.pool.GetTx(r)
	if err != nil {
		return nil, httputil.Error(http.StatusInternalServerError, err)
	}
	return r.WithContext(context.WithValue(r.Context(), txContextKey, tx)), nil
}

// ContextTx get the transaction on the request context
func (i Integrator) ContextTx(r *http.Request) (Tx, *httputil.HTTPError) {
	tx, ok := r.Context().Value(txContextKey).(Tx)
	if !ok {
		if tx == nil {
			return nil, httputil.Error(http.StatusInternalServerError, fmt.Errorf("cast to httpdb.Tx fail, was nil instead"))
		}
		return nil, httputil.Error(http.StatusInternalServerError, fmt.Errorf("cast to httpdb.Tx fail"))
	}
	return tx, nil
}

// TxRollbackRequestWrap wraps a httputil.RequestWrapFunc call Tx.FinalizeError() when f returns an error
func (i Integrator) TxRollbackRequestWrap(f httputil.RequestWrapFunc) httputil.RequestWrapFunc {
	return func(r *http.Request) (*http.Request, *httputil.HTTPError) {
		tx, httpErr := i.ContextTx(r)
		if httpErr != nil {
			return r, httpErr
		}

		req, err := f(r)
		if err == nil {
			return req, nil
		}

		if err := tx.FinalizeError(); err != nil {
			httputil.LogError(r, httputil.Error(http.StatusInternalServerError, err))
		}
		return req, err
	}
}

// TxFinalizeWrap wraps a httputil.ResponseJSONWrapFunc to:
// - call Tx.Finalize() if the request has no error
// - otherwise, call tx.FinalizeError()
func (i Integrator) TxFinalizeWrap(f httputil.ResponseJSONWrapFunc) httputil.ResponseJSONWrapFunc {
	return func(r *http.Request) (any, *httputil.HTTPError) {
		tx, httpErr := i.ContextTx(r)
		if httpErr != nil {
			return nil, httpErr
		}
		resp, httpErr := f(r)
		if httpErr != nil {
			_ = tx.FinalizeError() // only call on failure
			return resp, httpErr
		}
		if err := tx.Finalize(); err != nil {
			_ = tx.FinalizeError() // only call on failure
			return nil, httputil.Error(http.StatusInternalServerError, err)
		}
		return resp, nil
	}
}
