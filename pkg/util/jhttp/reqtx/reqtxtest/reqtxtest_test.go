package reqtxtest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type txMode int

type Txn struct{ name string }

func (t Txn) Finalize() error      { return nil }
func (t Txn) FinalizeError() error { return nil }

func newRequest() *http.Request { return httptest.NewRequest(http.MethodGet, "https://fake.com", nil) }

func TestPool_SetTxGetTx(t *testing.T) {
	expectedTx := Txn{name: "my_name"}
	pool := NewPool[Txn, txMode]()
	mode := txMode(1)
	txReq := pool.SetTx(t, newRequest(), expectedTx, mode)

	testCases := []struct {
		name string
		req  *http.Request
		mode txMode
		err  error
	}{
		{name: "normal", mode: mode},
		{name: "diff_request", req: newRequest(), mode: mode, err: errors.New("transaction with id, , does not exist")},
		{name: "diff_mode", mode: -9, err: errors.New("stored Tx mode (1) is not matching passed mode (-9)")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			req := tc.req
			if req == nil {
				req = txReq
			}

			tx, err := pool.GetTx(req, tc.mode)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			require.Equal(expectedTx, tx)
		})
	}
}
