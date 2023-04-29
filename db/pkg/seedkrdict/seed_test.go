package seedkrdict

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestSeed(t *testing.T) {
	testName := "TestSeed"
	testSeed(t, testName, func() error {
		return Seed(context.Background(), fixture.JoinTestData(testName))
	})
}

func TestSeedFile(t *testing.T) {
	testName := "TestSeedFile"
	testSeed(t, testName, func() error {
		return SeedFile(context.Background(), fixture.Read(t, "TestSeedFile.xml"))
	})
}

func testSeed(t *testing.T, testName string, f func() error) {
	require := require.New(t)
	ctx := context.Background()
	testdb.SetupTempDBT(t, testName)

	err := f()
	require.NoError(err)

	queries := db.New(db.DB())

	count, err := queries.TermsCount(ctx)
	require.NoError(err)
	require.Equal(int64(3), count)

	terms, err := queries.TermsPopular(ctx)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, terms))
}

func TestUnmarshallRscPath(t *testing.T) {
	require := require.New(t)
	testName := "TestUnmarshallRscPath"

	lexes, err := UnmarshallRscPath(fixture.JoinTestData(testName))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, lexes))
}

func TestUnmarshallRscXML(t *testing.T) {
	require := require.New(t)
	testName := "TestUnmarshallRscXML"

	lex, err := UnmarshallRscXML(fixture.Read(t, testName+".xml"))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, lex))
}

func TestIsNoTranslationsFoundError(t *testing.T) {
	require := require.New(t)
	require.True(IsNoTranslationsFoundError(&NoTranslationsFoundError{}))
	require.False(IsNoTranslationsFoundError(fmt.Errorf("test error")))
}

func TestLexicalEntry_CreateParams(t *testing.T) {
	require := require.New(t)
	testName := "TestLexicalEntry_CreateParams"

	lex, err := UnmarshallRscXML(fixture.Read(t, testName+".xml"))
	require.NoError(err)

	createParamsArray := []db.TermCreateParams{}
	for i, entry := range lex.LexicalEntries {
		createParams, err := entry.CreateParams(i + 1)
		require.NoError(err)
		createParamsArray = append(createParamsArray, createParams)
	}
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParamsArray))
}

func TestLexicalEntry_Term(t *testing.T) {
	require := require.New(t)
	testName := "TestLexicalEntry_Term"

	lex, err := UnmarshallRscXML(fixture.Read(t, testName+".xml"))
	require.NoError(err)

	terms := []dictionary.Term{}
	for _, entry := range lex.LexicalEntries {
		term, err := entry.Term()
		require.NoError(err)
		terms = append(terms, term)
	}
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, terms))
}

func TestFindGoodExample(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

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
