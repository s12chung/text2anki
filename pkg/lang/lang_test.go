package lang

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToCommonLevel(t *testing.T) {
	require := require.New(t)

	commonLevel, err := ToCommonLevel(1)
	require.NoError(err)
	require.Equal(CommonLevelRare, commonLevel)

	commonLevel, err = ToCommonLevel(0)
	require.NoError(err)
	require.Equal(CommonLevelUnique, commonLevel)

	commonLevel, err = ToCommonLevel(3)
	require.NoError(err)
	require.Equal(CommonLevelCommon, commonLevel)

	commonLevel, err = ToCommonLevel(-1)
	require.Equal(errors.New("common level not within range 0 to 3: -1"), err)
	require.Equal(CommonLevelUnique, commonLevel)

	commonLevel, err = ToCommonLevel(4)
	require.Equal(errors.New("common level not within range 0 to 3: 4"), err)
	require.Equal(CommonLevelUnique, commonLevel)
}

func TestPartOfSpeechTypes(t *testing.T) {
	require := require.New(t)
	require.Len(PartOfSpeechTypes(), PartOfSpeechCount)
}

func TestToPartOfSpeech(t *testing.T) {
	require := require.New(t)

	pos, err := ToPartOfSpeech(string(PartOfSpeechAdverb))
	require.NoError(err)
	require.Equal(PartOfSpeechAdverb, pos)

	pos, err = ToPartOfSpeech("")
	require.NoError(err)
	require.Equal(PartOfSpeechEmpty, pos)

	pos, err = ToPartOfSpeech("NOT A POS")
	require.Equal(errors.New("pos not matching lang.PartOfSpeech: NOT A POS"), err)
	require.Equal(PartOfSpeechEmpty, pos)
}
