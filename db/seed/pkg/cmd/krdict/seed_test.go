package krdict

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestUnmarshallXML(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

	require := require.New(t)

	lex, err := unmarshallXML(fixture.Read(t, "TestUnmarshallXML.xml"))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestUnmarshallXML.json", fixture.JSON(t, lex))
}

func TestFindGoodExample(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

	entry := findGoodExample(t)
	fmt.Println(string(fixture.JSON(t, entry)))
	fixture.CompareReadOrUpdate(t, "TestFindGoodExample.json", fixture.JSON(t, entry))
}

func findGoodExample(t *testing.T) *lexicalEntry {
	require := require.New(t)

	lexes, err := unmarshallRscPath()
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
