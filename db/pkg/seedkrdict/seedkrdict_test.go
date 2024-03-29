package seedkrdict

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func init() { testdb.MustSetup() }

const rscPath = "../../" + DefaultRscPath

func TestRscXMLPaths(t *testing.T) {
	require := require.New(t)
	t.Parallel()
	paths, err := RscXMLPaths(fixture.JoinTestData("TestRscXMLPaths"))

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
