package khaiii

import (
	"encoding/json"
	"testing"

	"github.com/s12chung/text2anki/pkg/test/fixture"
	"github.com/stretchr/testify/require"
)

const pathChange = "../../"

func newKhaiii(t *testing.T) *Khaiii {
	require := require.New(t)
	k, err := NewKhaiii(pathChange + DefaultDlPath)
	require.Nil(err)
	return k
}

func TestAnalyze(t *testing.T) {
	require := require.New(t)

	var err error
	k := newKhaiii(t)
	err = k.Open(pathChange + DefaultRscPath)
	require.Nil(err)
	defer func() {
		require.Nil(k.Close())
	}()

	var words []*Word
	words, err = k.Analyze("안녕! 반가워!")
	require.Nil(err)

	bytes, err := json.MarshalIndent(words, "", "  ")
	require.Nil(err)

	fixture.CompareReadOrUpdate(t, "analyze.json", bytes)
}

func TestVersion(t *testing.T) {
	require := require.New(t)
	k := newKhaiii(t)
	require.Equal("0.5", k.Version())
}
