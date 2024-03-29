package azure

import (
	"context"
	"strings"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/s12chung/text2anki/pkg/util/test/vcr"
)

func TestAzure_TextToSpeech(t *testing.T) {
	t.Parallel()

	synth := New(GetAPIKeyFromEnv(), EastUSRegion)
	t.Cleanup(setupVCR(t, "TestAzure_TextToSpeech", synth))

	require := require.New(t)
	ctx := context.Background()

	speech, err := synth.TextToSpeech(ctx, "안녕")
	require.NoError(err)
	mtype := mimetype.Detect(speech)
	require.Equal(".mp3", mtype.Extension())
	require.Equal("audio/mpeg", mtype.String())

	// use cache
	_, err = synth.TextToSpeech(ctx, "안녕")
	require.NoError(err)
}

func setupVCR(t *testing.T, testName string, hasClient vcr.HasClient) func() {
	return vcr.SetupVCR(t, fixture.JoinTestData(testName), hasClient, func(r *recorder.Recorder) {
		r.AddHook(func(i *cassette.Interaction) error {
			delete(i.Request.Headers, apiKeyHeader)
			delete(i.Request.Headers, tokenHeader)
			return nil
		}, recorder.AfterCaptureHook)
		r.AddHook(func(i *cassette.Interaction) error {
			if strings.Contains(i.Request.URL, "issueToken") {
				i.Response.Body = "REDACTED"
			}
			return nil
		}, recorder.BeforeSaveHook)
	})
}
