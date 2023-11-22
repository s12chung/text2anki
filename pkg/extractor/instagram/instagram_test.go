package instagram

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
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

var extractToDirSuffixes = []int{0, 1, 2, 9, 10, 11, 99, 100, 101, 110}

const extractToDirPrefix = "2023-11-21_10-42-44_UTC_"

func init() {
	args := make([]string, len(extractToDirSuffixes)+1)
	args[0] = "touch"
	for i, suffix := range extractToDirSuffixes {
		args[i+1] = extractToDirPrefix + strconv.Itoa(suffix) + extensions[0]
	}

	extractToDirArgs = func(login, id string) []string { return args }
}

func TestPost_ExtractToDir(t *testing.T) {
	cacheDir := path.Join(os.TempDir(), test.GenerateName("Instagram"))
	require.NoError(t, os.MkdirAll(cacheDir, ioutil.OwnerRWXGroupRX))

	testCases := []struct {
		name string
		url  string
		err  error
	}{
		{name: "broken", url: "https://waka.com", err: fmt.Errorf("url is not verified for instagram: https://waka.com")},
		{name: "fake", url: testURL},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			err := NewPost(tc.url).ExtractToDir(cacheDir)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)

			entries, err := os.ReadDir(cacheDir)
			require.NoError(err)
			entryNames := make([]string, len(entries))
			for i, entry := range entries {
				entryNames[i] = entry.Name()
			}

			filenames := make([]string, len(extractToDirSuffixes))
			for i, suffix := range extractToDirSuffixes {
				filenames[i] = extractToDirPrefix + fmt.Sprintf("%03d", suffix) + extensions[0]
			}
			require.Equal(filenames, entryNames)
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
