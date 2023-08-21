// Package extractor extracts data to create db.SourcePartMedia
package extractor

import (
	"fmt"
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

// Source represents a souce to extract from
type Source interface {
	Verify() bool
	ID() string
	ExtractToDir(cacheDir string) error
}

// Extractor extracts Source data given the Factory Source
type Extractor struct {
	cacheDir string
	factory  Factory
}

// NewExtractor returns a new Extractor
func NewExtractor(cacheDir string, factory Factory) Extractor {
	return Extractor{cacheDir: cacheDir, factory: factory}
}

// Extract extracts data given the Factory Source
func (e Extractor) Extract(s string) ([]db.SourcePartMediaFile, error) {
	source := e.factory.NewSource(s)
	if !source.Verify() {
		return nil, fmt.Errorf("string does not match factory source: %v", s)
	}

	hash := source.ID()
	cacheDir := filepath.Join(e.cacheDir, hash)
	if err := os.MkdirAll(cacheDir, ioutil.OwnerRWXGroupRX); err != nil {
		return nil, err
	}
	if err := source.ExtractToDir(cacheDir); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(entries))
	for _, file := range entries {
		if file.IsDir() {
			continue
		}
		for _, ext := range e.factory.Extensions() {
			if strings.HasSuffix(file.Name(), ext) {
				files = append(files, filepath.Join(cacheDir, file.Name()))
				break
			}
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files that match extensions extracted: %v", strings.Join(e.factory.Extensions(), ", "))
	}
	parts := make([]db.SourcePartMediaFile, len(files))
	for i, file := range files {
		f, err := os.Open(file) //nolint:gosec // needed
		if err != nil {
			return nil, err
		}
		parts[i] = db.SourcePartMediaFile{ImageFile: f}
	}
	return parts, nil
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
