package testdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/krdict"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestSeed_Generate(t *testing.T) {
	require := require.New(t)
	lexes, err := krdict.UnmarshallRscPath(fixture.JoinTestData("TestSeed_Generate"))
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
	fixture.CompareReadOrUpdate(t, "Seed.json", fixture.JSON(t, terms))
}
