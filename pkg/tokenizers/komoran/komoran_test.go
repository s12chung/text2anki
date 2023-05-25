package komoran

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testPort = 9002

func TestKomoran_Tokenize(t *testing.T) {
	test.CISkip(t, "can't run java environment in CI")
	t.Parallel()
	require := require.New(t)

	tokenizer := newKomoran(testPort)
	err := tokenizer.Setup()
	defer func() {
		require.NoError(tokenizer.CleanupAndWait())
	}()
	require.NoError(err)
	tokens, err := tokenizer.Tokenize("대한민국은 민주공화국이다.")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestKomoran_Tokenize.json", fixture.JSON(t, tokens))
}
