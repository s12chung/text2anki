package koreanbasic

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/v2/cassette"
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/test/vcr"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	dict := New(GetAPIKeyFromEnv())
	clean := setupVCR(t, "TestSearch", dict)
	defer clean()

	testSearch(t, dict, "가다", "search_expected.json")
}

func TestSearchFail(t *testing.T) {
	dict := New(GetAPIKeyFromEnv())
	clean := setupVCR(t, "TestSearchFail", dict)
	defer clean()

	testSearch(t, dict, "가다", "search_fail_expected.json")
}

func testSearch(t *testing.T, dict dictionary.Dicionary, text, expectedFile string) {
	require := require.New(t)

	terms, err := dict.Search(text)
	require.Nil(err)
	resultBytes, err := json.MarshalIndent(terms, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, expectedFile, resultBytes)
}

func setupVCR(t *testing.T, testName string, hasClient interface{}) func() {
	return vcr.SetupVCR(t, fixture.JoinTestData(testName), hasClient, func(r *recorder.Recorder) {
		r.AddFilter(func(i *cassette.Interaction) error {
			i.URL = cleanURL(i.URL)
			return nil
		})
		r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
			return r.Method == i.Method && cleanURL(r.URL.String()) == i.URL
		})
	})
}

func cleanURL(url string) string {
	return strings.Replace(url, "key="+GetAPIKeyFromEnv(), "key=REDACTED", 1)
}
