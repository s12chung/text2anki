package koreanbasic

import (
	"context"
	"net/http"
	"path"
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
	t.Parallel()

	tcs := []struct {
		name       string
		searchTerm string
		pos        lang.PartOfSpeech
	}{
		{name: "basic", searchTerm: "가다"},
		// PartOfSpeechOther will convert to PartOfSpeechEmpty
		{name: "basic_with_other", searchTerm: "가다", pos: lang.PartOfSpeechOther},
		{name: "empty", searchTerm: "안녕하세요"},
		{name: "auxiliary_verb", searchTerm: "가다", pos: lang.PartOfSpeechAuxiliaryVerb},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			t.Parallel()

			filenameName := tc.name
			if tc.name == "basic_with_other" {
				filenameName = "basic"
			}

			dict := New(GetAPIKeyFromEnv())
			clean := vcr.SetupVCR(t, fixture.JoinTestData(testName, filenameName), dict, func(r *recorder.Recorder) {
				r.AddHook(func(i *cassette.Interaction) error {
					i.Request.URL = cleanURL(i.Request.URL)
					return nil
				}, recorder.AfterCaptureHook)
				r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
					return r.Method == i.Method && cleanURL(r.URL.String()) == i.URL
				})
			})
			t.Cleanup(clean)

			terms, err := dict.Search(context.Background(), tc.searchTerm, tc.pos)
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, path.Join(testName, filenameName)+".json", fixture.JSON(t, terms))
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

func TestMergePosMap(t *testing.T) {
	require := require.New(t)

	require.Equal(lang.PartOfSpeechCount, len(mergePosMap))

	uniquePosMapValues := map[lang.PartOfSpeech]bool{}
	for _, v := range partOfSpeechMap {
		uniquePosMapValues[v] = true
	}
	uniquePosMapValues[lang.PartOfSpeechEmpty] = true

	uniquePosValues := map[lang.PartOfSpeech]bool{}
	for _, v := range mergePosMap {
		uniquePosValues[v] = true
	}
	require.Equal(uniquePosMapValues, uniquePosValues)
}
