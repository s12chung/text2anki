package xml

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestUnmarshall(t *testing.T) {
	require := require.New(t)

	bytes := fixture.Read(t, "TestSchema.json")
	node := &SchemaNode{}
	err := json.Unmarshal(bytes, node)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestSchema.json", fixture.JSON(t, node))
}

func TestSchema(t *testing.T) {
	require := require.New(t)

	node := NewSchemaNode()
	node.Attrs = Attrs{
		"test_merge": true,
	}

	node, err := Schema(fixture.Read(t, "TestSchema.xml"), node)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestSchema.json", fixture.JSON(t, node))
	testNodeManyCount(t, node)
}

func testNodeManyCount(t *testing.T, node *SchemaNode) {
	require := require.New(t)
	require.Zero(len(node.childrenMany))
	for _, child := range node.Children {
		testNodeManyCount(t, child)
	}
}
