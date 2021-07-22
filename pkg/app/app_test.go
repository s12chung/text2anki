package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"testing"

	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/dictionary/koreanbasic"
	"github.com/s12chung/text2anki/pkg/test/fixture"
)

func init() {
	if fixture.WillUpdate() {
		updateAnkiNotesTestdata()
	}
}

// updateAnkiNotesTestdata syncs the anki testdata to match with the korean dictionary ones, so that they match without
// being dependent on each other
func updateAnkiNotesTestdata() {
	sourcePath := path.Join("..", "dictionary", "koreanbasic", fixture.TestDataDir, "search.xml")
	//nolint:gosec // for tests
	sourceBytes, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Panic(fmt.Errorf("error while reading source fixture: %w", err))
	}
	terms, err := koreanbasic.SearchTerms(sourceBytes)
	if err != nil {
		log.Panic(fmt.Errorf("error while parsing source fixture: %w", err))
	}

	notes := make([]anki.Note, len(terms))
	for i, term := range terms {
		notes[i] = NewNoteFromTerm(term, 0)
	}
	fixtureBytes, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		log.Panic(fmt.Errorf("error while creating fixture: %w", err))
	}

	err = ioutil.WriteFile(path.Join("..", "anki", fixture.TestDataDir, "notes.json"), fixtureBytes, 0600)
	if err != nil {
		log.Panic(fmt.Errorf("error while writing fixture: %w", err))
	}
}

func TestFixtureCheck(t *testing.T) {
	if fixture.WillUpdate() {
		t.Fatal("Always fails if updating fixtures so that init() will run")
	}
}
