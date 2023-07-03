package db_test

import (
	"context"
	"encoding/json"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestTerm_StaticCopy(t *testing.T) {
	require := require.New(t)

	term := db.Term{}
	err := json.Unmarshal(fixture.Read(t, "TestToDBTerm.json"), &term)
	require.NoError(err)
	test.EmptyFieldsMatch(t, term)

	termCopy := term
	termCopy.ID = 0
	require.Equal(termCopy, term.StaticCopy())
}

func TestToDBTerm(t *testing.T) {
	require := require.New(t)
	testName := "TestToDBTerm"

	term := dictionary.Term{}
	err := json.Unmarshal(fixture.Read(t, testName+"Src.json"), &term)
	require.NoError(err)
	test.EmptyFieldsMatch(t, term)

	dbTerm, err := db.ToDBTerm(term, 1)
	require.NoError(err)
	test.EmptyFieldsMatch(t, dbTerm)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, dbTerm))
}

func TestTerm_DictionaryTerm(t *testing.T) {
	testName := "TestTerm_DictionaryTerm"
	tcs := []struct {
		name        string
		emptyFields []string
	}{
		{name: "Base"},
		{name: "EmptyVariants", emptyFields: []string{"Variants"}},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			term := db.Term{}
			err := json.Unmarshal(fixture.Read(t, path.Join(testName, tc.name+"Src.json")), &term)
			require.NoError(err)
			test.EmptyFieldsMatch(t, term, tc.emptyFields...)

			dictTerm, err := term.DictionaryTerm()
			require.NoError(err)
			test.EmptyFieldsMatch(t, dictTerm, tc.emptyFields...)

			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, dictTerm))
		})
	}
}

func TestTerm_CreateParams(t *testing.T) {
	require := require.New(t)
	testName := "TestTerm_CreateParams"

	term := db.Term{}
	err := json.Unmarshal(fixture.Read(t, "TestToDBTerm.json"), &term)
	require.NoError(err)
	test.EmptyFieldsMatch(t, term)

	createParams := term.CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestQueries_TermsSearchRaw(t *testing.T) {
	require := require.New(t)
	testName := "TestQueries_TermsSearch"
	ctx := context.Background()

	results, err := db.Qs().TermsSearchRaw(ctx, testdb.SearchTerm, lang.PartOfSpeechUnknown, testdb.SearchConfig)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, results))

	_, err = db.Qs().TermsSearchRaw(ctx, testdb.SearchTerm, lang.PartOfSpeechUnknown, db.TermsSearchConfig{
		PopWeight:    50,
		CommonWeight: 51,
	})
	require.Error(err)

	results, err = db.Qs().TermsSearchRaw(ctx, testdb.SearchTerm, testdb.SearchPOS, testdb.SearchConfig)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+"_Verb_"+".json", fixture.JSON(t, results))
}
