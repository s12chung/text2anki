package xz

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestRead(t *testing.T) {
	require := require.New(t)
	testName := "TestRead"

	bytes, err := Read(fixture.JoinTestData(testName + ".txt.xz"))
	require.NoError(err)
	require.Equal("xz_contents\n", string(bytes))
}
