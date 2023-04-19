package krdict

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestUnmarshallRscXML(t *testing.T) {
	require := require.New(t)

	lex, err := unmarshallRscXML(fixture.Read(t, "TestUnmarshallXML.xml"))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestUnmarshallXML.json", fixture.JSON(t, lex))
}

func TestTerm(t *testing.T) {
	require := require.New(t)

	lex, err := unmarshallRscXML(fixture.Read(t, "TestTerm.xml"))
	require.NoError(err)

	terms := []dictionary.Term{}
	for _, entry := range lex.LexicalEntries {
		term, err := entry.term()
		require.NoError(err)
		terms = append(terms, term)
	}
	fixture.CompareReadOrUpdate(t, "TestTerm.json", fixture.JSON(t, terms))
}

func TestFindGoodExample(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

	entry := findGoodExample(t)
	fmt.Println(string(fixture.JSON(t, entry)))
	fixture.CompareReadOrUpdate(t, "TestFindGoodExample.json", fixture.JSON(t, entry))
}

func findGoodExample(t *testing.T) *lexicalEntry {
	require := require.New(t)

	lexes, err := unmarshallRscPath(rscPath)
	require.NoError(err)
	for _, lex := range lexes {
		for _, entry := range lex.LexicalEntries {
			if goodExampleEntry(entry) && goodExampleSense(entry.Senses) && goodExampleWordForm(entry.WordForms) {
				return &entry
			}
		}
	}
	return nil
}

func goodExampleEntry(entry lexicalEntry) bool {
	return !(entry.Lemmas == nil || entry.RelatedForms == nil || entry.Senses == nil || entry.WordForms == nil)
}

func goodExampleSense(senses []sense) bool {
	for _, sense := range senses {
		if !(sense.Equivalents == nil || sense.Multimedias == nil || sense.SenseExamples == nil || sense.SenseRelations == nil) {
			return true
		}
	}
	return false
}

func goodExampleWordForm(wordForms []wordForm) bool {
	for _, wordForm := range wordForms {
		if !(wordForm.FormRepresentation.Feats == nil) {
			return true
		}
	}
	return false
}
