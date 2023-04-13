package krdict

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestUnmarshallXML(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

	require := require.New(t)

	lex, err := unmarshallXML(fixture.Read(t, "TestUnmarshallXML.xml"))
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, "TestUnmarshallXML.json", fixture.JSON(t, lex))
}

func TestFindGoodExample(t *testing.T) {
	test.CISkip(t, "rsc files not in CI")

	require := require.New(t)

	entry, err := findGoodExample()
	require.NoError(err)

	fmt.Println(string(fixture.JSON(t, entry)))
	fixture.CompareReadOrUpdate(t, "TestFindGoodExample.json", fixture.JSON(t, entry))
}
