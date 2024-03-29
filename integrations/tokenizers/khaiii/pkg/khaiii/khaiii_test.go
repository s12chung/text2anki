package khaiii

import (
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

func TestKhaiii_Analyze(t *testing.T) {
	test.CISkip(t, ciSkipMsg)

	require := require.New(t)

	var err error
	k := newKhaiii(t)
	require.NoError(k.Open(rscPath))
	defer func() { require.NoError(k.Close()) }()

	var words []Word
	words, err = k.Analyze("안녕! 반가워!")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestKhaiii_Analyze.json", fixture.JSON(t, words))
}

func TestKhaiii_Version(t *testing.T) {
	test.CISkip(t, ciSkipMsg)

	require := require.New(t)

	k := newKhaiii(t)
	require.Equal("0.5", k.Version())
}

func TestRscSHA(t *testing.T) {
	test.CISkip(t, ciSkipMsg)

	require := require.New(t)
	hashMap, err := fixture.SHA2Map(rscPath)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestRscSHA.json", fixture.JSON(t, hashMap))
}
