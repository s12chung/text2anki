package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/db/seed/pkg/cmd/krdict"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestQueries_TermsSearch(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	testName := "TestQueries_TermsSearch"
	testdb.SetupTempDBT(ctx, t, testName)

	err := krdict.Seed(context.Background(), fixture.JoinTestData(testName))
	require.NoError(err)

	queries := db.New(db.DB())
	results, err := queries.TermsSearch(ctx, "마음")
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestQueries_TermsSearch.json", fixture.JSON(t, results))
}
