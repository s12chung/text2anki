package komoran

import (
	"encoding/json"
	"testing"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestGetTokens(t *testing.T) {
	jarPath = "../../../tokenizers/dist/komoran"

	require := require.New(t)

	tokenizer := NewKomoran()
	err := tokenizer.Setup()
	defer func() {
		err = tokenizer.Cleanup()
		require.Nil(err)
	}()
	require.Nil(err)
	tokens, err := tokenizer.GetTokens("대한민국은 민주공화국이다.")
	require.Nil(err)

	bytes, err := json.MarshalIndent(tokens, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "get_tokens.json", bytes)
}
