package krdict

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestKrDict_Search(t *testing.T) {
	require := require.New(t)
	testName := "TestKrDict_Search"
	testdb.SetupTempDBT(t, testName)
	testdb.Seed(t)

	dict := New(db.DB())
	terms, err := dict.Search("마음")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, terms))
}
