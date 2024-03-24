package seedkrdict

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const xmlExt = ".xml"

func TestSeed(t *testing.T) {
	testName := "TestSeed"
	t.Parallel()

	testSeed(t, testName, func(tx db.Tx) error {
		return Seed(tx, fixture.JoinTestData(testName))
	})
}

func TestSeedFile(t *testing.T) {
	testName := "TestSeedFile"
	t.Parallel()

	testSeed(t, testName, func(tx db.Tx) error {
		return SeedFile(tx, fixture.Read(t, testName+xmlExt))
	})
}

func testSeed(t *testing.T, testName string, testFunc func(tx db.Tx) error) {
	require := require.New(t)

	txQs := testdb.TxQs(t, db.WriteOpts())
	require.NoError(txQs.TermsClearAll(txQs.Ctx()))

	require.NoError(testFunc(txQs))

	count, err := txQs.TermsCount(txQs.Ctx())
	require.NoError(err)
	require.Equal(int64(3), count)

	terms, err := txQs.TermsPopular(txQs.Ctx())
	require.NoError(err)
	fixture.CompareReadOrUpdateJSON(t, testName, terms)
}

func TestUnmarshallRscPath(t *testing.T) {
	require := require.New(t)
	testName := "TestUnmarshallRscPath"
	t.Parallel()

	lexes, err := UnmarshallRscPath(fixture.JoinTestData(testName))
	require.NoError(err)
	fixture.CompareReadOrUpdateJSON(t, testName, lexes)
}

func TestUnmarshallRscXML(t *testing.T) {
	require := require.New(t)
	testName := "TestUnmarshallRscXML"
	t.Parallel()

	lex, err := UnmarshallRscXML(fixture.Read(t, testName+xmlExt))
	require.NoError(err)
	fixture.CompareReadOrUpdateJSON(t, testName, lex)
}

func TestIsNoTranslationsFoundError(t *testing.T) {
	require := require.New(t)
	t.Parallel()
	require.True(IsNoTranslationsFoundError(NoTranslationsFoundError{}))
	require.False(IsNoTranslationsFoundError(fmt.Errorf("test error")))
}

func TestLexicalEntry_CreateParams(t *testing.T) {
	require := require.New(t)
	testName := "TestLexicalEntry_CreateParams"
	t.Parallel()

	lex, err := UnmarshallRscXML(fixture.Read(t, testName+xmlExt))
	require.NoError(err)

	createParamsArray := []db.TermCreateParams{}
	for i, entry := range lex.LexicalEntries {
		createParams, err := entry.CreateParams(i + 1)
		require.NoError(err)
		createParamsArray = append(createParamsArray, createParams)
	}
	fixture.CompareReadOrUpdateJSON(t, testName, createParamsArray)
}

func TestLexicalEntry_Term(t *testing.T) {
	require := require.New(t)
	testName := "TestLexicalEntry_Term"
	t.Parallel()

	lex, err := UnmarshallRscXML(fixture.Read(t, testName+xmlExt))
	require.NoError(err)

	terms := []dictionary.Term{}
	for _, entry := range lex.LexicalEntries {
		term, err := entry.Term()
		require.NoError(err)
		terms = append(terms, term)
	}
	fixture.CompareReadOrUpdateJSON(t, testName, terms)
}

func TestFindGoodExample(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")
	t.Parallel()

	entry := findGoodExample(t)
	// fmt.Println(string(fixture.JSON(t, entry)))
	fixture.CompareReadOrUpdate(t, "TestFindGoodExample.json", fixture.JSON(t, entry))
}

func findGoodExample(t *testing.T) *LexicalEntry {
	require := require.New(t)

	lexes, err := UnmarshallRscPath(rscPath)
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

func goodExampleEntry(entry LexicalEntry) bool {
	return !(entry.Lemmas == nil || entry.RelatedForms == nil || entry.Senses == nil || entry.WordForms == nil)
}

func goodExampleSense(senses []Sense) bool {
	for _, sense := range senses {
		if !(sense.Equivalents == nil || sense.Multimedias == nil || sense.SenseExamples == nil || sense.SenseRelations == nil) {
			return true
		}
	}
	return false
}

func goodExampleWordForm(wordForms []WordForm) bool {
	for _, wordForm := range wordForms {
		if !(wordForm.FormRepresentation.Feats == nil) {
			return true
		}
	}
	return false
}
