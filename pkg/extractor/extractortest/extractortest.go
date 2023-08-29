// Package extractortest is an implementation of extractor interfaces for tsting
package extractortest

import (
	"os"
	"path/filepath"

	"github.com/s12chung/text2anki/pkg/extractor"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

// Factory is the factory for  Source
type Factory struct{ fixturePath string }

// NewFactory returns a new factory
func NewFactory(name string) Factory {
	return Factory{fixturePath: fixture.JoinTestData(name)}
}

// NewSource returns a new Source
func (t Factory) NewSource(s string) extractor.Source {
	return Source{s: s, fixturePath: t.fixturePath}
}

// Extensions returns the extensions the extractor returns
func (t Factory) Extensions() []string { return []string{".jpg", ".png"} }

// Source represents a source to extract from
type Source struct {
	s           string
	fixturePath string
}

// VerifyString is the string to compare for Verify()
const VerifyString = "good_verify"

// SkipExtractString is the string to compare for Verify() and skipping ExtractToDir()
const SkipExtractString = "skip_extract"

// Verify returns true if the string matches VerifyString
func (t Source) Verify() bool { return t.s == VerifyString || t.s == SkipExtractString }

// ID returns a static id for the Source
func (t Source) ID() string { return t.s }

// ExtractToDir uses the fixture path as an extraction point to the cacheDir
func (t Source) ExtractToDir(cacheDir string) error {
	if t.s == SkipExtractString {
		return nil
	}
	err := filepath.Walk(t.fixturePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return ioutil.CopyFile(filepath.Join(cacheDir, info.Name()), path, ioutil.OwnerGroupR)
	})
	return err
}

// Info returns the info from the extraction
func (t Source) Info(_ string) (extractor.SourceInfo, error) {
	return extractor.SourceInfo{
		Name:      "extractortest",
		Reference: "extractortest.go source code",
	}, nil
}
