package text

import (
	"encoding/json"
	"testing"

	lingua "github.com/pemistahl/lingua-go"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestLanguagesMatch(t *testing.T) {
	require := require.New(t)
	require.Equal(int(Zulu)+1, len(lingua.AllLanguages()))
	require.Equal(int(Unknown), int(lingua.Unknown))
}

func TestTextsFromString(t *testing.T) {
	require := require.New(t)

	tcs := []struct {
		name string
	}{
		{name: "none"},
		{name: "simple_weave"},
		{name: "weave"},
		{name: "split"},
	}

	parser := NewParser(Korean, English)
	for _, tc := range tcs {
		s := string(fixture.Read(t, tc.name+".txt"))
		texts := parser.TextsFromString(s)

		bytes, err := json.MarshalIndent(texts, "", "  ")
		require.Nil(err)

		fixture.CompareReadOrUpdate(t, tc.name+".json", bytes)
	}
}
