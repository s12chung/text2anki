package koreanbasic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

const serachXML = "search.xml"

func init() {
	if fixture.WillUpdate() {
		koreanBasic := &KoreanBasic{apiKey: os.Getenv("KOREAN_BASIC_API_KEY")}
		bytes, err := koreanBasic.getSearch("가다")
		if err != nil {
			log.Panic(fmt.Errorf("error while getting fixture update: %w", err))
		}
		bytes = []byte("<!-- DO NOT EDIT. Generated in koreanbasic_test.go -->\n\n" + string(bytes))

		err = ioutil.WriteFile(fixture.JoinTestData(serachXML), bytes, 0600)
		if err != nil {
			log.Panic(fmt.Errorf("error while writing fixture: %w", err))
		}
	}
}

func TestParseSearch(t *testing.T) {
	require := require.New(t)

	channel, err := unmarshallSearch(fixture.Read(t, serachXML))
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(channel, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "search_expected.json", resultBytes)
}

func TestSearchTerms(t *testing.T) {
	require := require.New(t)

	terms, err := SearchTerms(fixture.Read(t, serachXML))
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(terms, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "search_terms.json", resultBytes)
}
