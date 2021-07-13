package koreanbasic

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestParseSearch(t *testing.T) {
	require := require.New(t)
	bytes, err := ioutil.ReadFile("testdata/search.xml")
	require.Nil(err)

	channel, err := parseSearch(bytes)
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(channel, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "testdata/search_expected.json", resultBytes)
}

func TestItemsToTerms(t *testing.T) {
	require := require.New(t)
	bytes, err := ioutil.ReadFile("testdata/search.xml")
	require.Nil(err)

	channel, err := parseSearch(bytes)
	require.Nil(err)
	terms := itemsToTerms(channel.Items)
	resultBytes, err := json.MarshalIndent(terms, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "testdata/search_items.json", resultBytes)
}
