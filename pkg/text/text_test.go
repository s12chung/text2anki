package text

import (
	"encoding/json"
	"testing"

	lingua "github.com/pemistahl/lingua-go"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestLanguagesMatch(t *testing.T) {
	require := require.New(t)
	require.Equal(int(Zulu)+1, len(lingua.AllLanguages()))
	require.Equal(int(Unknown), int(lingua.Unknown))
}

func TestTextsFromString(t *testing.T) {
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			s := string(fixture.Read(t, tc.name+".txt"))
			texts, err := parser.TextsFromString(s)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.Nil(err)

			bytes, err := json.MarshalIndent(texts, "", "  ")
			require.Nil(err)

			fixture.CompareReadOrUpdate(t, tc.name+".json", bytes)
		})
	}
}
