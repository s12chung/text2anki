// Package vcr contains helpers for go-vcr
package vcr

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

// HasClient objects has a client to set for VCR
type HasClient interface {
	SetClient(c *http.Client)
}

// SetupVCR setups up the VCR recorder
func SetupVCR(t *testing.T, cassetteName string, hasClient any, setupRecorder func(r *recorder.Recorder)) func() {
	require := require.New(t)
	h, ok := hasClient.(HasClient)
	require.True(ok, "should implement HasClient")

	r, err := recorder.New(cassetteName)
	require.NoError(err)

	if setupRecorder != nil {
		setupRecorder(r)
	}
	h.SetClient(&http.Client{Transport: r})

	return func() {
		require.NoError(r.Stop())
	}
}
