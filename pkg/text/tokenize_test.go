//nolint:dupword
package text

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
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

	fixture.CompareReadOrUpdate(t, "TestTokenizeTexts.json", fixture.JSON(t, tokenizedTexts))
}
