package db_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestToDBTerm(t *testing.T) {
	require := require.New(t)
	testName := "TestToDBTerm"

	term := dictionary.Term{}
	err := json.Unmarshal(fixture.Read(t, testName+"Src.json"), &term)
	require.NoError(err)

	dbTerm, err := db.ToDBTerm(term, 1)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, dbTerm))
}

func TestTerm_DictionaryTerm(t *testing.T) {
	require := require.New(t)
	testName := "TestTerm_DictionaryTerm"

	term := db.Term{}
	err := json.Unmarshal(fixture.Read(t, testName+"Src.json"), &term)
	require.NoError(err)

	dbTerm, err := term.DictionaryTerm()
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, dbTerm))
}

func TestTerm_CreateParams(t *testing.T) {
	require := require.New(t)
	testName := "TestTerm_CreateParams"

	term := db.Term{}
	err := json.Unmarshal(fixture.Read(t, "TestToDBTerm.json"), &term)
	require.NoError(err)

	createParams := term.CreateParams()
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestQueries_TermsSearch(t *testing.T) {
	require := require.New(t)
	testName := "TestQueries_TermsSearch"
	testdb.SetupTempDBT(t, testName)
	testdb.Seed(t)

	results, err := db.Qs().TermsSearch(context.Background(), "마음", db.TermsSearchConfig{
		PopLog:       20,
		PopWeight:    40,
		CommonWeight: 40,
		LenLog:       2,
	})
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, results))

	_, err = db.Qs().TermsSearch(context.Background(), "마음", db.TermsSearchConfig{
		PopWeight:    50,
		CommonWeight: 51,
	})
	require.Error(err)
}
