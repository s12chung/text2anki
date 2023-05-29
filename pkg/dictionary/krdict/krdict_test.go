package krdict

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

type MustSetupAndSeed struct{}

func TestMain(m *testing.M) {
	testdb.MustSetupAndSeed(MustSetupAndSeed{})
	os.Exit(m.Run())
}

func TestKrDict_Search(t *testing.T) {
	require := require.New(t)
	testName := "TestKrDict_Search"
	dict := New(db.DB())

	// PartOfSpeechOther will convert to PartOfSpeechEmpty
	for _, pos := range []lang.PartOfSpeech{lang.PartOfSpeechEmpty, lang.PartOfSpeechOther} {
		terms, err := dict.Search("마음", pos)
		require.NoError(err)
		fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, terms))
	}
}

func TestMergePosMap(t *testing.T) {
	require := require.New(t)

	require.Equal(lang.PartOfSpeechCount, len(mergePosMap))

	posMap := seedkrdict.PartOfSpeechMap()
	uniquePosMapValues := map[lang.PartOfSpeech]bool{}
	for _, v := range posMap {
		uniquePosMapValues[v] = true
	}
	uniquePosMapValues[lang.PartOfSpeechEmpty] = true

	uniquePosValues := map[lang.PartOfSpeech]bool{}
	for _, v := range mergePosMap {
		uniquePosValues[v] = true
	}
	require.Equal(uniquePosMapValues, uniquePosValues)
}
