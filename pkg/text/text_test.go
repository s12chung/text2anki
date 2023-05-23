package text

import (
	"testing"

	"github.com/pemistahl/lingua-go"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestLanguagesMatch(t *testing.T) {
	require := require.New(t)
	require.Equal(int(Zulu)+1, len(lingua.AllLanguages()))
	require.Equal(int(Unknown), int(lingua.Unknown))
}

func TestParser_TextsFromString(t *testing.T) {
	testNamePath := "TestParser_TextsFromString/"
	tcs := []struct {
		name string
		err  error
	}{
		{name: "none"},
		{name: "simple_weave"},
		{name: "weave"},
		{name: "split"},
		{name: "split_1_line"},
		{name: "split_extra_text", err: errExtraTextLine},
		{name: "split_extra_translation", err: errExtraTranslationLine},
		{name: "split_1_line_extra_translation", err: errExtraTranslationLine},
	}

	parser := NewParser(Korean, English)
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			s := string(fixture.Read(t, testNamePath+tc.name+".txt"))
			texts, err := parser.TextsFromString(s)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.Nil(err)

			fixture.CompareReadOrUpdate(t, testNamePath+tc.name+".json", fixture.JSON(t, texts))
		})
	}
}
