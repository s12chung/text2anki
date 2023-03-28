package koreanbasic

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/v2/cassette"
	"github.com/dnaeon/go-vcr/v2/recorder"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/s12chung/text2anki/pkg/util/test/vcr"
)

func TestSearch(t *testing.T) {
	dict := New(GetAPIKeyFromEnv())
	clean := vcr.SetupVCR(t, fixture.JoinTestData("TestSearch"), dict, func(r *recorder.Recorder) {
		r.AddFilter(func(i *cassette.Interaction) error {
			i.URL = cleanURL(i.URL)
			return nil
		})
		r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
			return r.Method == i.Method && cleanURL(r.URL.String()) == i.URL
		})
	})
	defer clean()

	tcs := []struct {
		searchTerm string
		expected   string
	}{
		{searchTerm: "가다", expected: "search_expected.json"},
		{searchTerm: "안녕하세요", expected: "search_empty_expected.json"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.expected, func(t *testing.T) {
			require := require.New(t)

			terms, err := dict.Search(tc.searchTerm)
			require.NoError(err)
			resultBytes, err := json.MarshalIndent(terms, "", "  ")
			require.NoError(err)

			fixture.CompareReadOrUpdate(t, tc.expected, resultBytes)
		})
	}
}

func cleanURL(url string) string {
	return strings.Replace(url, "key="+GetAPIKeyFromEnv(), "key=REDACTED", 1)
}
