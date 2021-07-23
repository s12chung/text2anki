package koreanbasic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

const searchXML = "search.xml"
const searchXMLFail = "search_fail.xml"

func init() {
	if fixture.WillUpdateAPI() {
		koreanBasic := &KoreanBasic{apiKey: GetAPIKeyFromEnv()}
		files := map[string]string{
			searchXML:     "가다",
			searchXMLFail: "안녕하세요",
		}
		for filename, term := range files {
			bytes, err := koreanBasic.getSearch(term)
			if err != nil {
				log.Panic(fmt.Errorf("error while getting fixture update: %w", err))
			}
			bytes = []byte("<!-- DO NOT EDIT. Generated in koreanbasic_test.go -->\n\n" + string(bytes))

			err = ioutil.WriteFile(fixture.JoinTestData(filename), bytes, 0600)
			if err != nil {
				log.Panic(fmt.Errorf("error while writing fixture: %w", err))
			}
		}
	}
}

func TestParseSearch(t *testing.T) {
	require := require.New(t)

	tcs := []struct {
		fixture  string
		expected string
	}{
		{fixture: searchXML, expected: "search_expected.json"},
		{fixture: searchXMLFail, expected: "search_fail_expected.json"},
	}

	for _, tc := range tcs {
		channel, err := unmarshallSearch(fixture.Read(t, tc.fixture))
		require.Nil(err)
		resultBytes, err := json.MarshalIndent(channel, "", "  ")
		require.Nil(err)

		fixture.CompareReadOrUpdate(t, tc.expected, resultBytes)
	}
}

func TestSearchTerms(t *testing.T) {
	require := require.New(t)

	terms, err := SearchTerms(fixture.Read(t, searchXML))
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(terms, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "search_terms.json", resultBytes)
}
