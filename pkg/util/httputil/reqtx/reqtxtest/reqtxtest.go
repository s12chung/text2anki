// Package reqtxtest contains testing struct for reqtx
package reqtxtest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/httputil/reqtx"
)

// Pool is a pool that maps transactions to an ID stored as idHeader in request headers
type Pool struct{ idMap map[string]reqtx.Tx }

// NewPool returns a new Pool
func NewPool() Pool { return Pool{idMap: map[string]reqtx.Tx{}} }

// SetTxT is SetTx with a *testing.T shorthand
func (p Pool) SetTxT(t *testing.T, r *http.Request, tx reqtx.Tx) *http.Request {
	require := require.New(t)
	require.NoError(p.SetTx(r, tx))
	return r
}

// SetTx maps the transaction with a new ID, the ID is set to the idHeader in the request header
func (p Pool) SetTx(r *http.Request, tx reqtx.Tx) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	idString := id.String()

	p.idMap[idString] = tx
	r.Header.Set(idHeader, idString)
	return nil
}

const idHeader = "X-Request-ID"

// GetTx returns the transaction given the id stored in idHeader
func (p Pool) GetTx(r *http.Request) (reqtx.Tx, error) {
	id := r.Header.Get(idHeader)
	tx, exists := p.idMap[id]
	if !exists {
		return nil, fmt.Errorf("transaction with id, %v, does not exist", id)
	}
	return tx, nil
}
