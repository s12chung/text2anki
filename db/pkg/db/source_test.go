package db

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

var textTokenizer = TextTokenizer{
	Parser:       text.NewParser(text.Korean, text.English),
	Tokenizer:    tokenizers.NewSplitTokenizer(),
	CleanSpeaker: true,
}

func TestTextTokenizer_TokenizeTextsFromString(t *testing.T) {
	testNamePath := "TestTextTokenizer_TokenizeTextsFromString/"

	testCases := []struct {
		name string
	}{
		{name: "split"},
		{name: "weave"},
		{name: "speaker_split"},
		{name: "speaker_weave"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			s := string(fixture.Read(t, testNamePath+tc.name+".txt"))
			tokenizedTexts, err := textTokenizer.TokenizeTextsFromString(s)
			require.NoError(err)

			nonSpeaker := strings.TrimPrefix(tc.name, "speaker_")
			fixture.CompareReadOrUpdate(t, testNamePath+nonSpeaker+".json", fixture.JSON(t, tokenizedTexts))
		})
	}
}

func TestTextTokenizer_TokenizeTexts(t *testing.T) {
	require := require.New(t)

	//nolint:dupword
	texts := []text.Text{
		{Text: "super mario go"},
		{Text: "waka waka"},
		{Text: "all this to learn korean"},
	}

	tokenizedTexts, err := textTokenizer.TokenizeTexts(texts)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestTextTokenizer_TokenizeTexts.json", fixture.JSON(t, tokenizedTexts))
}
