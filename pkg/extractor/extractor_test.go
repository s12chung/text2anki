package extractor_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/extractor/extractortest"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestSourceExtraction_InfoFile(t *testing.T) {
	require := require.New(t)

	info := db.PrePartInfo{
		Name:      "extractor_test_name",
		Reference: "https://extactor-test.com",
	}
	f, err := extractor.SourceExtraction{Info: info}.InfoFile()
	require.NoError(err)
	b, err := io.ReadAll(f)
	require.NoError(err)

	fileInfo := db.PrePartInfo{}
	require.NoError(json.Unmarshal(b, &fileInfo))
	require.Equal(info, fileInfo)
}

func TestExtractor_Extract(t *testing.T) {
	testName := "TestExtractor_Extract"
	cacheDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, cacheDir)

	testCases := []struct {
		name string
		s    string
		err  error
	}{
		{name: "basic", s: extractortest.VerifyString},
		{name: "skip_extract", s: extractortest.SkipExtractString, err: fmt.Errorf("no filenames that match extensions extracted: .jpg, .png")},
		{name: "no_verify", s: "fail", err: fmt.Errorf("string does not match factory source: fail")},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			source, err := extractor.NewExtractor(cacheDir, extractortest.NewFactory(testName)).Extract(tc.s)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+"_info.json"), fixture.JSON(t, source.Info))

			partMap := map[string]string{}
			for _, part := range source.Parts {
				require.Nil(part.AudioFile)

				file := part.ImageFile
				info, err := file.Stat()
				require.NoError(err)
				bytes, err := io.ReadAll(file)
				require.NoError(err)
				partMap[info.Name()] = string(bytes)
			}
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+"_parts.json"), fixture.JSON(t, partMap))
		})
	}
}

func TestVerify(t *testing.T) {
	testName := "TestVerify"
	cacheDir := path.Join(os.TempDir(), test.GenerateName(testName))
	test.MkdirAll(t, cacheDir)

	extractorMap := extractor.Map{
		"test": extractor.NewExtractor(cacheDir, extractortest.NewFactory(testName)),
	}

	testCases := []struct {
		name     string
		s        string
		expected string
	}{
		{name: "basic", s: extractortest.VerifyString, expected: "test"},
		{name: "no_verify", s: "fail"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, extractor.Verify(tc.s, extractorMap))
		})
	}
}
