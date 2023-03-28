package komoran

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testPort = 9002

func TestGetTokens(t *testing.T) {
	if test.IsCI() {
		t.Skip("can't run java environment in CI")
	}

	jarPath = "../../../" + jarPath
	require := require.New(t)

	tokenizer := new(testPort)
	err := tokenizer.Setup()
	defer func() {
		require.NoError(tokenizer.CleanupAndWait())
	}()
	require.NoError(err)
	tokens, err := tokenizer.Tokenize("대한민국은 민주공화국이다.")
	require.NoError(err)

	bytes, err := json.MarshalIndent(tokens, "", "  ")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "get_tokens.json", bytes)
}
