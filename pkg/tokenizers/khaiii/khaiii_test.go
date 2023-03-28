package khaiii

import (
	"encoding/json"
	"testing"

	"github.com/s12chung/text2anki/pkg/test"
	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestGetTokens(t *testing.T) {
	if test.IsCI() {
		t.Skip("can't run java environment in CI")
	}

	binPath = "../../../" + binPath
	require := require.New(t)

	tokenizer := new()
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
