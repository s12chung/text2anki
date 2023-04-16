package krdict

import (
	"testing"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
	"github.com/stretchr/testify/require"
)

func init() {
	rscPath = "../../../../" + rscPath
}

func TestRscXMLPaths(t *testing.T) {
	oldRscPath := rscPath
	rscPath = fixture.JoinTestData("TestRscXMLPaths")
	defer func() {
		rscPath = oldRscPath
	}()

	require := require.New(t)
	paths, err := RscXMLPaths()

	require.NoError(err)
	require.Equal([]string{
		"testdata/TestRscXMLPaths/1125431_5000.xml",
		"testdata/TestRscXMLPaths/1125431_10000.xml",
		"testdata/TestRscXMLPaths/1125431_15000.xml",
		"testdata/TestRscXMLPaths/1125431_20000.xml",
		"testdata/TestRscXMLPaths/1125431_25000.xml",
		"testdata/TestRscXMLPaths/1125431_30000.xml",
		"testdata/TestRscXMLPaths/1125431_35000.xml",
		"testdata/TestRscXMLPaths/1125431_40000.xml",
		"testdata/TestRscXMLPaths/1125431_45000.xml",
		"testdata/TestRscXMLPaths/1125431_50000.xml",
		"testdata/TestRscXMLPaths/1125431_51960.xml",
	}, paths)
}
