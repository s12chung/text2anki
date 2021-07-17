package koreanbasic

import (
	"encoding/json"
	"testing"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestParseSearch(t *testing.T) {
	require := require.New(t)

	channel, err := unmarshallSearch(fixture.Read(t, "search.xml"))
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(channel, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "search_expected.json", resultBytes)
}

func TestSearchTerms(t *testing.T) {
	require := require.New(t)

	terms, err := SearchTerms(fixture.Read(t, "search.xml"))
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(terms, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "search_terms.json", resultBytes)
}
