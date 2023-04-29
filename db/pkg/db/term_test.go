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
	ctx := context.Background()
	testName := "TestQueries_TermsSearch"
	testdb.SetupTempDBT(t, testName)
	testdb.Seed(t)

	queries := db.New(db.DB())
	results, err := queries.TermsSearch(ctx, "마음")
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, results))
}
