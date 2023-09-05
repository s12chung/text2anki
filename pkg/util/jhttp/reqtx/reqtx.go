// Package reqtx manages a database transaction per request
package reqtx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

const txContextKey jhttp.ContextKey = "reqtx.TxContext"

// Tx represents a transaction
type Tx interface {
	Finalize() error
	FinalizeError() error
}

// Pool gets transactions to handle
type Pool[T Tx] interface {
	GetTx(r *http.Request) (T, error)
}

// Integrator integrates Pool to the request via the Router
type Integrator[T Tx] struct{ pool Pool[T] }

// NewIntegrator returns a new Integrator
func NewIntegrator[T Tx](pool Pool[T]) Integrator[T] { return Integrator[T]{pool: pool} }

// SetTxContext sets the transaction on the request context
func (i Integrator[T]) SetTxContext(r *http.Request) (*http.Request, *jhttp.HTTPError) {
	tx, err := i.pool.GetTx(r)
	if err != nil {
		return nil, jhttp.Error(http.StatusInternalServerError, err)
	}
	return r.WithContext(context.WithValue(r.Context(), txContextKey, tx)), nil
}

// ContextTx get the transaction on the request context
func ContextTx[T Tx](r *http.Request) (T, *jhttp.HTTPError) {
	value := r.Context().Value(txContextKey)
	tx, ok := value.(T)
	if !ok {
		var empty T
		if value == nil {
			return empty, jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to httpdb.Tx fail, was nil instead"))
		}
		return empty, jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to httpdb.Tx fail"))
	}
	return tx, nil
}

// TxRollbackRequestWrap wraps a jhttp.RequestWrapFunc call Tx.FinalizeError() when f returns an error
func TxRollbackRequestWrap(f jhttp.RequestWrapFunc) jhttp.RequestWrapFunc {
	return func(r *http.Request) (*http.Request, *jhttp.HTTPError) {
		tx, httpErr := ContextTx[Tx](r)
		if httpErr != nil {
			return r, httpErr
		}

		req, err := f(r)
		if err == nil {
			return req, nil
		}

		if err := tx.FinalizeError(); err != nil {
			jhttp.LogError(r, jhttp.Error(http.StatusInternalServerError, err))
		}
		return req, err
	}
}

// TxFinalizeWrap wraps a jhttp.ResponseJSONWrapFunc to:
// - call Tx.Finalize() if the request has no error
// - otherwise, call tx.FinalizeError()
func TxFinalizeWrap(f jhttp.ResponseJSONWrapFunc) jhttp.ResponseJSONWrapFunc {
	return func(r *http.Request) (any, *jhttp.HTTPError) {
		tx, httpErr := ContextTx[Tx](r)
		if httpErr != nil {
			return nil, httpErr
		}
		model, httpErr := f(r)
		if httpErr != nil {
			_ = tx.FinalizeError() // only call on failure
			return model, httpErr
		}
		if err := tx.Finalize(); err != nil {
			_ = tx.FinalizeError() // only call on failure
			return nil, jhttp.Error(http.StatusInternalServerError, err)
		}
		return model, nil
	}
}
