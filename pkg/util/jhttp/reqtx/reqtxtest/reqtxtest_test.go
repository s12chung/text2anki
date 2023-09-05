package reqtxtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type txn struct{ name string }

func (t *txn) Finalize() error      { return nil }
func (t *txn) FinalizeError() error { return nil }

func newRequest() *http.Request { return httptest.NewRequest(http.MethodGet, "https://fake.com", nil) }

func TestPool_SetTxGetTx(t *testing.T) {
	require := require.New(t)

	expectedTx := &txn{name: "my_name"}
	pool := NewPool()
	req := pool.SetTxT(t, newRequest(), expectedTx)

	tx, err := pool.GetTx(req)
	require.NoError(err)
	require.Equal(expectedTx, tx)
}
