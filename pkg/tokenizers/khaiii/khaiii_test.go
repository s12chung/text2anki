package khaiii

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testPort = 9001

func TestGetTokens(t *testing.T) {
	test.CISkip(t, "can't run C environment in CI")

	binPath = "../../../" + binPath
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
