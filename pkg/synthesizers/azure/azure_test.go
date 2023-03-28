package azure

import (
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/v2/cassette"
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/s12chung/text2anki/pkg/util/test/vcr"
)

func TestTextToSpeech(t *testing.T) {
	synth := New(GetAPIKeyFromEnv(), EastUSRegion)
	clean := setupVCR(t, "TestTextToSpeech", synth)
	defer clean()

	require := require.New(t)

	speech, err := synth.TextToSpeech("안녕")
	require.NoError(err)
	mtype := mimetype.Detect(speech)
	require.Equal(".mp3", mtype.Extension())
	require.Equal("audio/mpeg", mtype.String())

	// use cache
	_, err = synth.TextToSpeech("안녕")
	require.NoError(err)
}

func setupVCR(t *testing.T, testName string, hasClient interface{}) func() {
	return vcr.SetupVCR(t, fixture.JoinTestData(testName), hasClient, func(r *recorder.Recorder) {
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
	})
}
