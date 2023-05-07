package koreanbasic

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/s12chung/text2anki/pkg/util/test/vcr"
)

func TestKoreanBasic_Search(t *testing.T) {
	testName := "TestKoreanBasic_Search"
	dict := New(GetAPIKeyFromEnv())
	clean := vcr.SetupVCR(t, fixture.JoinTestData(testName), dict, func(r *recorder.Recorder) {
		r.AddHook(func(i *cassette.Interaction) error {
			i.Request.URL = cleanURL(i.Request.URL)
			return nil
		}, recorder.AfterCaptureHook)
		r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
			return r.Method == i.Method && cleanURL(r.URL.String()) == i.URL
		})
	})
	defer clean()

	tcs := []struct {
		searchTerm string
		pos        lang.PartOfSpeech
		expected   string
	}{
		{searchTerm: "가다", expected: testName + "/expected.json"},
		{searchTerm: "안녕하세요", expected: testName + "/empty_expected.json"},
		{searchTerm: "가다", pos: lang.PartOfSpeechAuxiliaryVerb, expected: testName + "/pos_expected.json"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.expected, func(t *testing.T) {
			require := require.New(t)

			terms, err := dict.Search(tc.searchTerm, tc.pos)
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, tc.expected, fixture.JSON(t, terms))
		})
	}
}

func cleanURL(url string) string {
	return strings.Replace(url, "key="+GetAPIKeyFromEnv(), "key=REDACTED", 1)
}

func TestPartOfSpeechToAPIIntMatch(t *testing.T) {
	require := require.New(t)

	for k := range partOfSpeechToAPIInt {
		_, exists := partOfSpeechMap[k]
		require.True(exists, "For key, %v", k)
	}
}
