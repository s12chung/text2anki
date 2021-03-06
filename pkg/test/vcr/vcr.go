// Package vcr contains helpers for go-vcr
package vcr

import (
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/stretchr/testify/require"
)

// HasClient objects has a client to set for VCR
type HasClient interface {
	SetClient(c *http.Client)
}

// SetupVCR setups up the VCR recorder
func SetupVCR(t *testing.T, cassetteName string, hasClient interface{}, setupRecorder func(r *recorder.Recorder)) func() {
	require := require.New(t)
	h, ok := hasClient.(HasClient)
	require.True(ok, "should implement HasClient")

	r, err := recorder.New(cassetteName)
	require.Nil(err)

	if setupRecorder != nil {
		setupRecorder(r)
	}
	h.SetClient(&http.Client{Transport: r})

	return func() {
		require.Nil(r.Stop())
	}
}
