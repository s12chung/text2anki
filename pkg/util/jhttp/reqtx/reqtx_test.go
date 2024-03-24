package reqtx

import (
	"context"
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

	ctx  context.Context //nolint:containedctx // for testing
	mode txMode
}

func (t *txn) Finalize() error {
	if t.mode == finalizeErrorMode {
		return fmt.Errorf("test: finalizeErrorMode")
	}

	t.finalizeCount++
	return nil
}
func (t *txn) FinalizeError() error {
	t.finalizeErrorCount++
	return nil
}

type txMode int

const okMode txMode = 0
const finalizeErrorMode txMode = -1

type pool struct{ txn *txn }

func (p *pool) GetTx(r *http.Request, mode txMode) (Tx, error) {
	if mode != okMode && mode != finalizeErrorMode {
		return nil, fmt.Errorf("test: !okMode")
	}
	if p.txn == nil {
		p.txn = &txn{ctx: r.Context(), mode: mode}
	}
	return p.txn, nil
}

func newIntegrator() Integrator[Tx, txMode] { return NewIntegrator[Tx, txMode](&pool{}) }
func newRequest() *http.Request             { return httptest.NewRequest(http.MethodGet, "https://fake.com", nil) }

func TestIntegrator_SetTxModeContext(t *testing.T) {
	require := require.New(t)

	integrator := newIntegrator()
	req := newRequest()

	expectedMode := txMode(1)
	req = integrator.SetTxModeContext(req, expectedMode)

	mode, err := TxMode[txMode](req)
	require.Nil(err)
	require.Equal(expectedMode, mode)
}

func TestTxMode(t *testing.T) {
	testCases := []struct {
		name         string
		mode         any
		expectedMode txMode
		err          error
	}{
		{name: "normal", expectedMode: 1},
		{name: "nil", mode: nil},
		{name: "string", mode: "fail", err: jhttp.Error(http.StatusInternalServerError, fmt.Errorf("cast to reqtx.txMode fail"))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			m := tc.mode
			if tc.expectedMode != 0 {
				m = tc.expectedMode
			}
			req := newRequest()
			req = req.WithContext(context.WithValue(req.Context(), txOptsContextKey, m))

			mode, err := TxMode[txMode](req)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.Nil(err)
			require.Equal(tc.expectedMode, mode)
		})
	}
}

func TestIntegrator_ResponseWrap(t *testing.T) {
	testCases := []struct {
		name string
		req  *http.Request
		mode txMode

		reqErr *jhttp.HTTPError

		err                   *jhttp.HTTPError
		errMode               txMode
		errFinalizeErrorCount int
	}{
		{name: "normal"},
		{name: "req_error", reqErr: jhttp.Error(http.StatusBadRequest, fmt.Errorf("test_induced"))},
		{name: "bad_mode", mode: -9, err: jhttp.Error(http.StatusInternalServerError, fmt.Errorf("test: !okMode"))},
		{name: "finalize_fail", mode: finalizeErrorMode,
			errMode:               finalizeErrorMode,
			errFinalizeErrorCount: 1,
			err:                   jhttp.Error(http.StatusInternalServerError, fmt.Errorf("test: finalizeErrorMode"))},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			integrator := newIntegrator()
			req := tc.req
			if req == nil {
				req = newRequest()
			}
			req = integrator.SetTxModeContext(req, tc.mode)

			model, httpErr := integrator.ResponseWrap(func(_ *http.Request, _ Tx) (any, *jhttp.HTTPError) {
				return nil, tc.reqErr
			})(req)

			tx, err := integrator.GetTx(req, okMode)
			require.NoError(err)

			require.Equal(nil, model)
			if tc.err != nil {
				require.Equal(tc.err, httpErr)
				require.Equal(&txn{ctx: req.Context(), mode: tc.errMode, finalizeErrorCount: tc.errFinalizeErrorCount}, tx)
				return
			}
			if tc.reqErr != nil {
				require.Equal(tc.reqErr, httpErr)
				require.Equal(&txn{ctx: req.Context(), finalizeErrorCount: 1}, tx)
				return
			}
			require.Nil(httpErr)
			require.Equal(&txn{ctx: req.Context(), finalizeCount: 1}, tx)
		})
	}
}
