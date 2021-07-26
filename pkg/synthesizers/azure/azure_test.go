package azure

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/v2/cassette"
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/require"
)

func TestTextToSpeech(t *testing.T) {
	require := require.New(t)
	r, clean := setupVCR(t)
	defer clean()

	synth := NewAzure(GetAPIKeyFromEnv(), EastUSRegion)

	azure, ok := synth.(*Azure)
	require.True(ok)
	azure.client = &http.Client{Transport: r}

	speech, err := synth.TextToSpeech("안녕")
	require.Nil(err)
	mtype := mimetype.Detect(speech)
	require.Equal(".mp3", mtype.Extension())
	require.Equal("audio/mpeg", mtype.String())
}

func setupVCR(t *testing.T) (*recorder.Recorder, func()) {
	require := require.New(t)

	r, err := recorder.New("testdata/TestTextToSpeech")
	require.Nil(err)
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, apiKeyHeader)
		delete(i.Request.Headers, tokenHeader)
		return nil
	})
	r.AddSaveFilter(func(i *cassette.Interaction) error {
		if strings.Contains(i.URL, "issueToken") {
			i.Response.Body = "REDACTED"
		}
		return nil
	})
	return r, func() {
		require.Nil(r.Stop())
	}
}
