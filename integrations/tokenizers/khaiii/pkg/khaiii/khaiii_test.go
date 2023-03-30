package khaiii

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const pathChange = "../../"
const rscPath = pathChange + DefaultRscPath
const ciSkipMsg = "can't run C environment in CI"

func newKhaiii(t *testing.T) *Khaiii {
	require := require.New(t)
	k, err := NewKhaiii(pathChange + DefaultDlPath)
	require.NoError(err)
	return k
}

func TestAnalyze(t *testing.T) {
	test.CISkip(t, ciSkipMsg)

	require := require.New(t)

	var err error
	k := newKhaiii(t)
	require.NoError(k.Open(rscPath))
	defer func() {
		require.NoError(k.Close())
	}()

	var words []Word
	words, err = k.Analyze("안녕! 반가워!")
	require.NoError(err)

	bytes, err := json.MarshalIndent(words, "", "  ")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestAnalyze.json", bytes)
}

func TestVersion(t *testing.T) {
	test.CISkip(t, ciSkipMsg)

	require := require.New(t)

	k := newKhaiii(t)
	require.Equal("0.5", k.Version())
}

func TestRsc(t *testing.T) {
	test.CISkip(t, ciSkipMsg)

	require := require.New(t)
	hashMap, err := fixture.SHA2Map(rscPath)
	require.NoError(err)

	bytes, err := json.MarshalIndent(hashMap, "", "  ")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestRsc.json", bytes)
}
