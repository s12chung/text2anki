package text

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/s12chung/text2anki/pkg/tokenizers"
)

func TestTokenizeTexts(t *testing.T) {
	require := require.New(t)

	texts := []Text{
		{Text: "super mario go"},
		{Text: "waka waka"},
		{Text: "all this to learn korean"},
	}
	tokenizedTexts, err := TokenizeTexts(tokenizers.NewSplitTokenizer(), texts)
	require.Nil(err)
	bytes, err := json.MarshalIndent(tokenizedTexts, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "TestTokenizeTexts.json", bytes)
}
