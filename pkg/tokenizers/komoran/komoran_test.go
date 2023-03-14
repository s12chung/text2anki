package komoran

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

	jarPath = "../../../" + jarPath
	require := require.New(t)

	tokenizer := new()
	err := tokenizer.Setup()
	defer func() {
		require.Nil(tokenizer.server.StopAndWait())
	}()
	require.Nil(err)
	tokens, err := tokenizer.Tokenize("대한민국은 민주공화국이다.")
	require.Nil(err)

	bytes, err := json.MarshalIndent(tokens, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "get_tokens.json", bytes)
}
