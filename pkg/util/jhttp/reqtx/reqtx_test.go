package reqtx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

type txn struct {
	finalizeCount      int
	finalizeErrorCount int
	id                 string
}

func (t *txn) Finalize() error {
	t.finalizeCount++
	return nil
}
func (t *txn) FinalizeError() error {
	t.finalizeErrorCount++
	return nil
}

const txNameKey = "Tx-ID"
const txID = "12345"

type pool struct{}

func (p pool) GetTx(r *http.Request) (Tx, error) {
	id := r.Header.Get(txNameKey)
	if id == "" {
		return nil, nil //nolint:nilnil // for testing
	}
	return &txn{id: id}, nil
}

func newIntegrator() Integrator { return NewIntegrator(pool{}) }
func newRequest() *http.Request { return httptest.NewRequest(http.MethodGet, "https://fake.com", nil) }
func withTx(r *http.Request) *http.Request {
	r.Header.Set(txNameKey, txID)
	return r
}

func TestIntegrator_SetTxContext(t *testing.T) {
	require := require.New(t)

	integrator := newIntegrator()

	req, err := integrator.SetTxContext(withTx(newRequest()))
	require.Nil(err)

	tx, err := ContextTx(req)
	require.Nil(err)
	require.Equal(&txn{id: txID}, tx)
}

func TestTxRollbackRequestWrap(t *testing.T) {
	testCases := []struct {
		name    string
		request *http.Request
		err     *jhttp.HTTPError
		reqErr  *jhttp.HTTPError
	}{
		{name: "normal", request: withTx(newRequest())},
		{name: "no_id", request: newRequest(), err: jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to httpdb.Tx fail, was nil instead"))},
		{name: "req_error", request: withTx(newRequest()), reqErr: jhttp.Error(http.StatusBadRequest, fmt.Errorf("test_induced"))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			integrator := newIntegrator()
			req, err := integrator.SetTxContext(tc.request)
			require.Nil(err)

			finalReq, err := TxRollbackRequestWrap(func(r *http.Request) (*http.Request, *jhttp.HTTPError) {
				return r, tc.reqErr
			})(req)

			tx, ctxErr := ContextTx(finalReq)

			require.Equal(req, finalReq)
			if tc.err != nil {
				require.Equal(tc.err, err)
				require.Equal(tc.err, ctxErr)
				require.Nil(tx)
				return
			}
			if tc.reqErr != nil {
				require.Equal(tc.reqErr, err)
				require.Nil(ctxErr)
				require.Equal(&txn{id: txID, finalizeErrorCount: 1}, tx)
				return
			}
			require.Nil(err)
			require.Nil(ctxErr)
			require.Equal(&txn{id: txID}, tx)
		})
	}
}

func TestTxFinalizeWrap(t *testing.T) {
	testCases := []struct {
		name    string
		request *http.Request
		err     *jhttp.HTTPError
		reqErr  *jhttp.HTTPError
	}{
		{name: "normal", request: withTx(newRequest())},
		{name: "no_id", request: newRequest(), err: jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to httpdb.Tx fail, was nil instead"))},
		{name: "req_error", request: withTx(newRequest()), reqErr: jhttp.Error(http.StatusBadRequest, fmt.Errorf("test_induced"))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			integrator := newIntegrator()
			req, err := integrator.SetTxContext(tc.request)
			require.Nil(err)

			model, err := TxFinalizeWrap(func(r *http.Request) (any, *jhttp.HTTPError) {
				return nil, tc.reqErr
			})(req)

			tx, ctxErr := ContextTx(req)

			require.Equal(nil, model)
			if tc.err != nil {
				require.Equal(tc.err, err)
				require.Equal(tc.err, ctxErr)
				require.Nil(tx)
				return
			}
			if tc.reqErr != nil {
				require.Equal(tc.reqErr, err)
				require.Nil(ctxErr)
				require.Equal(&txn{id: txID, finalizeErrorCount: 1}, tx)
				return
			}
			require.Nil(err)
			require.Nil(ctxErr)
			require.Equal(&txn{id: txID, finalizeCount: 1}, tx)
		})
	}
}
