package text

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestParseSubtitles(t *testing.T) {
	tcs := []struct {
		name string
	}{
		{name: "match"},
		{name: "simple"},
		{name: "complex"},
	}
	const fixtureName = "TestParseSubtitles.json"

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			texts, err := ParseSubtitles(
				fixture.JoinTestData(tc.name+"_source.vtt"),
				fixture.JoinTestData(tc.name+"_translation.vtt"),
			)
			require.Nil(err)
			bytes, err := json.MarshalIndent(texts, "", "  ")
			require.Nil(err)

			if tc.name == "match" {
				fixture.SafeUpdate(t, fixtureName, bytes)
			}
			fixture.CompareRead(t, fixtureName, bytes)
		})
	}
}
