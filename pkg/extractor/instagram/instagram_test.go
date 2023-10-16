package instagram

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testLogin = "instalogin"
const testProtocol = "https://"
const testHash = "nF1jE8cVT8O"
const testPath = pathPrefix + testHash + "/"
const testQueryParams = "?igshid=nF1jE3HOcXWZjiBqjZQayKs"
const testURL = testProtocol + hostname + testPath + testQueryParams

func NewPost(url string) extractor.Source { return NewFactory(testLogin).NewSource(url) }

func TestSource_Verify(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		ok   bool
	}{
		{name: "basic", url: testURL, ok: true},
		{name: "no_path", url: testProtocol + hostname, ok: false},
		{name: "wrong_host", url: "https://waka.com" + pathPrefix + "abcd", ok: false},
		{name: "invalid", url: "waka.com/p/abcd", ok: false},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.ok, NewPost(tc.url).Verify())
		})
	}
}

func TestSource_ID(t *testing.T) {
	testCases := []struct {
		name   string
		url    string
		result string
	}{
		{name: "basic", url: testURL, result: testHash},
		{name: "no_slash", url: testProtocol + hostname + strings.TrimRight(testPath, "/"), result: testHash},
		{name: "broken", url: "https://waka.com"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.result, NewPost(tc.url).ID())
		})
	}
}

func TestPost_ExtractToDir(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		err  error
	}{
		{name: "broken", url: "https://waka.com", err: fmt.Errorf("url is not vertified for instagram: https://waka.com")},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.err, NewPost(tc.url).ExtractToDir(""))
		})
	}
}

func TestPost_Info(t *testing.T) {
	require := require.New(t)
	testName := "TestPost_Info"

	info, err := NewPost("https://testpostinfo.com").Info(fixture.JoinTestData(testName))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, filepath.Join(testName, "postInfo.json"), fixture.JSON(t, info))
}
