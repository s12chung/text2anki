// Package reqtxtest contains testing struct for reqtx
package reqtxtest

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/jhttp/reqtx"
)

// Pool is a pool that maps transactions to an ID stored as idHeader in request headers
type Pool[T reqtx.Tx] struct {
	idMap map[string]T
	mutex *sync.RWMutex
}

// NewPool returns a new Pool
func NewPool[T reqtx.Tx]() Pool[T] { return Pool[T]{idMap: map[string]T{}, mutex: &sync.RWMutex{}} }

// SetTxT is SetTx with a *testing.T shorthand
func (p Pool[T]) SetTxT(t *testing.T, r *http.Request, tx T) *http.Request {
	require := require.New(t)
	require.NoError(p.SetTx(r, tx))
	return r
}

// SetTx maps the transaction with a new ID, the ID is set to the idHeader in the request header
func (p Pool[T]) SetTx(r *http.Request, tx T) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	idString := id.String()

	p.mutex.Lock()
	p.idMap[idString] = tx
	p.mutex.Unlock()

	r.Header.Set(idHeader, idString)
	return nil
}

const idHeader = "X-Request-ID"

// GetTx returns the transaction given the id stored in idHeader
func (p Pool[T]) GetTx(r *http.Request) (T, error) {
	id := r.Header.Get(idHeader)

	p.mutex.RLock()
	tx, exists := p.idMap[id]
	p.mutex.RUnlock()

	if !exists {
		var empty T
		return empty, fmt.Errorf("transaction with id, %v, does not exist", id)
	}
	return tx, nil
}
