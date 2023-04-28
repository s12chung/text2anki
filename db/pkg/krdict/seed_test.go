package krdict

import (
	"context"
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
	ctx := context.Background()
	testSeed(ctx, t, testName, func() error {
		return Seed(ctx, fixture.JoinTestData(testName))
	})
}

func TestSeedFile(t *testing.T) {
	testName := "TestSeedFile"
	ctx := context.Background()
	testSeed(ctx, t, testName, func() error {
		return SeedFile(ctx, fixture.Read(t, "TestSeedFile.xml"))
	})
}

func testSeed(ctx context.Context, t *testing.T, testName string, f func() error) {
	require := require.New(t)
	testdb.SetupTempDBT(ctx, t, testName)

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

func TestUnmarshallRscXML(t *testing.T) {
	require := require.New(t)

	lex, err := unmarshallRscXML(fixture.Read(t, "TestUnmarshallXML.xml"))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestUnmarshallXML.json", fixture.JSON(t, lex))
}

func TestLexicalEntry_term(t *testing.T) {
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
	// fmt.Println(string(fixture.JSON(t, entry)))
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
