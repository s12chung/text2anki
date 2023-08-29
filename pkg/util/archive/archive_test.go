package archive

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestXZReader(t *testing.T) {
	require := require.New(t)
	testName := "TestXZReader"

	bytes, err := XZBytes(fixture.JoinTestData(testName + ".txt.xz"))
	require.NoError(err)
	require.Equal("xz_contents\n", string(bytes))
}
