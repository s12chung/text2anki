package papago

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/s12chung/text2anki/pkg/util/test/vcr"
)

func TestPapago_Translate(t *testing.T) {
	t.Parallel()

	translator := New(GetClientCredentialsFromEnv())
	t.Cleanup(setupVCR(t, "TestPapago_Translate", translator))

	require := require.New(t)
	ctx := context.Background()

	text := "한 송이 꽃을 피우려 작은 두 눈에 얼마나 많은 비가 내렸을까"
	translation, err := translator.Translate(ctx, text)
	require.NoError(err)
	require.Equal("How much rain would have fallen in my two small eyes to bloom", translation)

	// use cache
	_, err = translator.Translate(ctx, text)
	require.NoError(err)
}

func setupVCR(t *testing.T, testName string, hasClient vcr.HasClient) func() {
	return vcr.SetupVCR(t, fixture.JoinTestData(testName), hasClient, func(r *recorder.Recorder) {
		r.AddHook(func(i *cassette.Interaction) error {
			delete(i.Request.Headers, clientIDHeader)
			delete(i.Request.Headers, clientSecretHeader)
			return nil
		}, recorder.AfterCaptureHook)
	})
}
