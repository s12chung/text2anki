package krdict

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestSchema(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

	require := require.New(t)

	hashMap, err := fixture.SHA2Map(rscPath)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestRscSchemaSHA.json", fixture.JSON(t, hashMap))

	node, err := RscSchema()
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestRscSchema.json", fixture.JSON(t, node))
}
