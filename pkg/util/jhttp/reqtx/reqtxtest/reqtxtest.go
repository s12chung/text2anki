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

type txValue[T reqtx.Tx, Mode ~int] struct {
	Tx   T
	Mode Mode
}

// Pool is a pool that maps transactions to an ID stored as idHeader in request headers
type Pool[T reqtx.Tx, Mode ~int] struct {
	idMap map[string]txValue[T, Mode]
	mutex *sync.RWMutex
}

// NewPool returns a new Pool
func NewPool[T reqtx.Tx, Mode ~int]() Pool[T, Mode] {
	return Pool[T, Mode]{idMap: map[string]txValue[T, Mode]{}, mutex: &sync.RWMutex{}}
}

// SetTx is SetTx with a *testing.T shorthand
func (p Pool[T, Mode]) SetTx(t *testing.T, r *http.Request, tx T, mode Mode) *http.Request {
	require := require.New(t)

	id, err := uuid.NewV7()
	require.NoError(err)
	idString := id.String()

	p.mutex.Lock()
	p.idMap[idString] = txValue[T, Mode]{Tx: tx, Mode: mode}
	p.mutex.Unlock()

	r.Header.Set(idHeader, idString)
	return r
}

const idHeader = "X-Request-ID"

// GetTx returns the transaction given the id stored in idHeader
func (p Pool[T, Mode]) GetTx(r *http.Request, mode Mode) (T, error) {
	id := r.Header.Get(idHeader)

	p.mutex.RLock()
	val, exists := p.idMap[id]
	p.mutex.RUnlock()

	var empty T
	if !exists {
		return empty, fmt.Errorf("transaction with id, %v, does not exist", id)
	}
	if val.Mode != mode {
		return empty, fmt.Errorf("stored Tx mode (%v) is not matching passed mode (%v)", val.Mode, mode)
	}
	return val.Tx, nil
}
