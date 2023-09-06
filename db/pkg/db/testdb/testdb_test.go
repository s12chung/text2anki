package testdb

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func writeModelFile(t *testing.T, filename string, fileBytes []byte) {
	require := require.New(t)
	require.NoError(os.WriteFile(filepath.Join(modelsDir, filename), fileBytes, ioutil.OwnerRWGroupR))
}

func TestGen___NotesSeed(t *testing.T) {
	require := require.New(t)
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	notes, err := Notes().Models()
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
	writeModelFile(t, Terms().Filename(), fixture.JSON(t, terms))
}

func TestGen___SourceStructuredsSeed(t *testing.T) {
	testName := "TestGen___SourceStructuredsSeed"
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	require := require.New(t)
	ctx := context.Background()

	tokenizer := db.TextTokenizer{Parser: config.Parser(), Tokenizer: config.Tokenizer(ctx, config.TokenizerKhaiii)}
	require.NoError(tokenizer.Setup(ctx))
	defer func() { require.NoError(tokenizer.Cleanup()) }()

	filepaths := allFilePaths(t, fixture.JoinTestData(testName))
	sources := make([]db.SourceStructured, len(filepaths))
	for i, fp := range filepaths {
		split := strings.Split(string(test.Read(t, fp)), "===")
		if len(split) == 1 {
			split = append(split, "")
		}
		tokenizedTexts, err := tokenizer.TokenizedTexts(ctx, split[0], split[1])
		require.NoError(err)

		reference := path.Base(fp)
		name := reference[0 : len(reference)-len(filepath.Ext(reference))]
		sources[i] = db.SourceStructured{Name: name, Reference: reference, Parts: []db.SourcePart{{TokenizedTexts: tokenizedTexts}}}
		test.EmptyFieldsMatch(t, sources[i], "ID", "UpdatedAt", "CreatedAt")
	}
	writeModelFile(t, SourceStructureds().Filename(), fixture.JSON(t, sources))
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
