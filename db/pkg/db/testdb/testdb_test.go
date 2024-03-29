package testdb

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/api/config"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestGen___NotesSeed(t *testing.T) {
	require := require.New(t)

	notes, err := Notes().Models()
	require.NoError(err)
	for _, note := range notes {
		emptyFields := []string{"ID", "Downloaded", "UpdatedAt", "CreatedAt"}
		if note.CommonLevel == 0 {
			emptyFields = append(emptyFields, "CommonLevel")
		}
		test.EmptyFieldsMatch(t, note, emptyFields...)
	}
}

func TestGen___TermsSeed(t *testing.T) {
	testName := "TestGen___TermsSeed"
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

			emptyFields := []string{"Variants"}
			if dbTerm.CommonLevel == 0 {
				emptyFields = append(emptyFields, "CommonLevel")
			}
			test.EmptyFieldsMatch(t, dbTerm, emptyFields...)

			terms = append(terms, dbTerm)
		}
		basePopularity += len(lex.LexicalEntries)
	}
	compareReadOrUpdate(t, Terms().Filename(), fixture.JSON(t, terms))
}

func TestGen___SourceStructuredsSeed(t *testing.T) {
	test.CISkip(t, "can't run tokenizer.TokenizedTexts in CI")

	testName := "TestGen___SourceStructuredsSeed"
	require := require.New(t)
	ctx := context.Background()

	tokenizer := db.TextTokenizer{Parser: config.Parser(), Tokenizer: config.Tokenizer(ctx, config.TokenizerKhaiii, plog)}
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

		if i == len(filepaths)-1 {
			sources[i].Parts[0].Media = &db.SourcePartMedia{ImageKey: SourcePartMediaImageKey}
			test.EmptyFieldsMatch(t, sources[i].Parts[0])
		} else {
			test.EmptyFieldsMatch(t, sources[i].Parts[0], "Media")
		}
	}
	compareReadOrUpdate(t, SourceStructureds().Filename(), fixture.JSON(t, sources))
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

func compareReadOrUpdate(t *testing.T, filename string, resultBytes []byte) {
	require := require.New(t)
	p := path.Join(modelsDir, filename)

	if fixture.WillUpdate() {
		assert := assert.New(t)
		assert.NoError(os.WriteFile(p, resultBytes, ioutil.OwnerRWGroupR))
		assert.Fail(fixture.UpdateFailMessage)
	}
	expected, err := os.ReadFile(p) //nolint:gosec // for tests
	require.NoError(err)
	require.Equal(strings.TrimSpace(string(expected)), strings.TrimSpace(string(resultBytes)))
}
