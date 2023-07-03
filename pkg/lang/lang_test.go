package lang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPartOfSpeechTypes(t *testing.T) {
	require := require.New(t)
	got := PartOfSpeechTypes()
	require.Equal(PartOfSpeechCount, len(got))
}
