package testdb

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/api"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestGen___TermsSeed(t *testing.T) {
	testName := "TestGen___TermsSeed"
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	require := require.New(t)
	lexes, err := seedkrdict.UnmarshallRscPath(fixture.JoinTestData(testName))
	require.NoError(err)

	var terms []db.Term

	basePopularity := 1
	for _, lex := range lexes {
		for i, entry := range lex.LexicalEntries {
			term, err := entry.Term()
			require.NoError(err)
			dbTerm, err := db.ToDBTerm(term, basePopularity+i)
			require.NoError(err)
			terms = append(terms, dbTerm)
		}
		basePopularity += len(lex.LexicalEntries)
	}
	fixture.Update(t, "TermsSeed.json", fixture.JSON(t, terms))
}

func TestGen___SourceSerializedsSeed(t *testing.T) {
	testName := "TestGen___SourceSerializedsSeed"
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	require := require.New(t)
	require.NoError(api.DefaultRoutes.Setup())
	defer func() {
		require.NoError(api.DefaultRoutes.Cleanup())
	}()

	filepaths := allFilePaths(t, fixture.JoinTestData(testName))
	sources := make([]db.SourceSerialized, len(filepaths))
	for i, fp := range filepaths {
		tokenizedTexts, err := api.DefaultRoutes.TextTokenizer.TokenizeTextsFromString(string(test.Read(t, fp)))
		require.NoError(err)
		sources[i] = db.SourceSerialized{Name: path.Base(fp), TokenizedTexts: tokenizedTexts}
	}
	fixture.Update(t, sourceSerializedsSeedFilename, fixture.JSON(t, sources))
}

func allFilePaths(t *testing.T, p string) []string {
	require := require.New(t)

	paths := []string{}
	files, err := os.ReadDir(p)
	require.NoError(err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		paths = append(paths, path.Join(p, path.Join(file.Name())))
	}
	return paths
}
