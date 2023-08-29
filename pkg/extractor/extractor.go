// Package extractor extracts data to create db.SourcePartMedia
package extractor

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

// Factory is a factory for Source
type Factory interface {
	NewSource(s string) Source
	Extensions() []string
}

// SourceInfo contains Source related info from the extractor
type SourceInfo struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

// Source represents a source to extract from
type Source interface {
	Verify() bool
	ID() string
	ExtractToDir(cacheDir string) error
	Info(cacheDir string) (SourceInfo, error)
}

// Extractor extracts Source data given the Factory Source
type Extractor struct {
	cacheDir string
	factory  Factory
}

// SourceExtraction is all the data extracted
type SourceExtraction struct {
	Info  SourceInfo               `json:"info"`
	Parts []db.SourcePartMediaFile `json:"parts"`
}

// InfoFile returns a file that contains .Info as json
func (s SourceExtraction) InfoFile() (fs.File, error) {
	f, err := os.CreateTemp("", "text2anki-InfoFile-*.json")
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(s.Info)
	if err != nil {
		return nil, err
	}
	if _, err := f.Write(b); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return os.Open(f.Name())
}

// NewExtractor returns a new Extractor
func NewExtractor(cacheDir string, factory Factory) Extractor {
	return Extractor{cacheDir: cacheDir, factory: factory}
}

// Extract extracts data given the Factory Source
func (e Extractor) Extract(s string) (SourceExtraction, error) {
	source := e.factory.NewSource(s)
	if !source.Verify() {
		return SourceExtraction{}, fmt.Errorf("string does not match factory source: %v", s)
	}

	hash := source.ID()
	cacheDir := filepath.Join(e.cacheDir, hash)
	if err := os.MkdirAll(cacheDir, ioutil.OwnerRWXGroupRX); err != nil {
		return SourceExtraction{}, err
	}
	if err := source.ExtractToDir(cacheDir); err != nil {
		return SourceExtraction{}, err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return SourceExtraction{}, err
	}
	filenames := filenamesWithExtensions(entries, e.factory.Extensions())
	if len(filenames) == 0 {
		return SourceExtraction{}, fmt.Errorf("no filenames that match extensions extracted: %v", strings.Join(e.factory.Extensions(), ", "))
	}

	parts := make([]db.SourcePartMediaFile, len(filenames))
	for i, filename := range filenames {
		f, err := os.Open(filepath.Join(cacheDir, filename)) //nolint:gosec // needed
		if err != nil {
			return SourceExtraction{}, err
		}
		parts[i] = db.SourcePartMediaFile{ImageFile: f}
	}
	info, err := source.Info(cacheDir)
	if err != nil {
		return SourceExtraction{}, err
	}
	return SourceExtraction{Info: info, Parts: parts}, nil
}

func filenamesWithExtensions(entries []os.DirEntry, extensions []string) []string {
	filenames := make([]string, 0, len(entries))
	for _, file := range entries {
		if file.IsDir() {
			continue
		}
		for _, ext := range extensions {
			if strings.HasSuffix(file.Name(), ext) {
				filenames = append(filenames, file.Name())
				break
			}
		}
	}
	return filenames
}

// Map is a map of extractor name to Extractor
type Map map[string]Extractor

// Verify returns the key matching the Extractor that Verify() == true
func Verify(s string, extractorMap Map) string {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "\n") {
		return ""
	}
	for k, extractor := range extractorMap {
		if extractor.factory.NewSource(s).Verify() {
			return k
		}
	}
	return ""
}
