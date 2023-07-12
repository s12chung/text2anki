package testdb

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/models"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/api"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const modelsPath = "models/modeldata"

func writeModelFile(t *testing.T, filename string, fileBytes []byte) {
	require := require.New(t)
	require.NoError(os.WriteFile(filepath.Join(modelsPath, filename), fileBytes, ioutil.OwnerRWGroupR))
}

func TestGen___NotesSeed(t *testing.T) {
	require := require.New(t)
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	notes, err := models.Notes()
	require.NoError(err)
	for _, note := range notes {
		emptyFields := []string{"ID", "Downloaded"}
		if note.CommonLevel == 0 {
			emptyFields = append(emptyFields, "CommonLevel")
		}
		test.EmptyFieldsMatch(t, note, emptyFields...)
	}
}

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

			emptyFields := []string{"Variants"}
			if dbTerm.CommonLevel == 0 {
				emptyFields = append(emptyFields, "CommonLevel")
			}
			test.EmptyFieldsMatch(t, dbTerm, emptyFields...)
		}
		basePopularity += len(lex.LexicalEntries)
	}
	writeModelFile(t, models.TermsSeedFilename, fixture.JSON(t, terms))
}

func TestGen___SourceSerializedsSeed(t *testing.T) {
	testName := "TestGen___SourceSerializedsSeed"
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	require := require.New(t)

	routes := api.NewRoutes(api.Config{})
	require.NoError(routes.Setup())
	defer func() {
		require.NoError(routes.Cleanup())
	}()

	filepaths := allFilePaths(t, fixture.JoinTestData(testName))
	sources := make([]db.SourceSerialized, len(filepaths))
	for i, fp := range filepaths {
		split := strings.Split(string(test.Read(t, fp)), "===")
		if len(split) == 1 {
			split = append(split, "")
		}
		tokenizedTexts, err := routes.TextTokenizer.TokenizedTexts(split[0], split[1])
		require.NoError(err)
		sources[i] = db.SourceSerialized{Name: path.Base(fp), Parts: []db.SourcePart{{TokenizedTexts: tokenizedTexts}}}
		test.EmptyFieldsMatch(t, sources[i], "ID", "UpdatedAt", "CreatedAt")
	}
	writeModelFile(t, models.SourceSerializedsSeedFilename, fixture.JSON(t, sources))
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
