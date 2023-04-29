package testdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestGen___Seed(t *testing.T) {
	if !fixture.WillUpdate() {
		t.Skip("TestGen___ test generates fixtures")
	}
	require := require.New(t)
	lexes, err := seedkrdict.UnmarshallRscPath(fixture.JoinTestData("TestGen___Seed"))
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
	fixture.Update(t, "Seed.json", fixture.JSON(t, terms))
}
