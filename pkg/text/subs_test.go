package text

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestParseSubtitles(t *testing.T) {
	testName := "TestParseSubtitles"
	tcs := []struct {
		name string
	}{
		{name: "match"},
		{name: "simple"},
		{name: "complex"},
	}
	fixtureName := testName + ".json"

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			texts, err := ParseSubtitles(
				fixture.JoinTestData(testName, tc.name+"_source.vtt"),
				fixture.JoinTestData(testName, tc.name+"_translation.vtt"),
			)
			require.NoError(err)

			bytes := fixture.JSON(t, texts)
			if tc.name == "match" {
				fixture.SafeUpdate(t, fixtureName, bytes)
			}
			fixture.CompareRead(t, fixtureName, bytes)
		})
	}
}
