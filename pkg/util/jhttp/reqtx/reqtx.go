// Package reqtx manages a database transaction per request
package reqtx

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

// Tx represents a transaction
type Tx interface {
	Finalize() error
	FinalizeError() error
}

// Pool gets transactions to handle
type Pool[T Tx, Mode ~int] interface {
	GetTx(r *http.Request, mode Mode) (T, error)
}

// Integrator integrates Pool to the request via the Router
type Integrator[T Tx, Mode ~int] struct{ Pool[T, Mode] }

// NewIntegrator returns a new Integrator
func NewIntegrator[T Tx, Mode ~int](pool Pool[T, Mode]) Integrator[T, Mode] {
	return Integrator[T, Mode]{Pool: pool}
}

const txOptsContextKey jhttp.ContextKey = "reqtx.TxOptsContext"

// SetTxModeContext sets the transaction mode on the request context
func (i Integrator[T, Mode]) SetTxModeContext(r *http.Request, mode Mode) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), txOptsContextKey, mode))
}

// TxMode returns the transaction mode from the request
func TxMode[Mode ~int](r *http.Request) (Mode, *jhttp.HTTPError) {
	value := r.Context().Value(txOptsContextKey)
	if value == nil {
		return 0, nil
	}
	mode, ok := value.(Mode)
	if !ok {
		return 0, jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to %v fail", reflect.TypeOf(mode).String()))
	}
	return mode, nil
}

// ResponseHandler is the handler for responses with transactions
type ResponseHandler[T Tx] func(r *http.Request, tx T) (any, *jhttp.HTTPError)

// ResponseWrap wraps a ResponseHandler to:
// - get the Tx and pass to the ResponseHandler
// - call Tx.Finalize() if the request has no error
// - otherwise, call tx.FinalizeError()
func (i Integrator[T, Mode]) ResponseWrap(f ResponseHandler[T]) jhttp.ResponseHandler {
	return func(r *http.Request) (any, *jhttp.HTTPError) {
		mode, httpErr := TxMode[Mode](r)
		if httpErr != nil {
			return nil, httpErr
		}
		tx, err := i.GetTx(r, mode)
		if err != nil {
			return nil, jhttp.Error(http.StatusInternalServerError, err)
		}

		model, httpErr := f(r, tx)
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
