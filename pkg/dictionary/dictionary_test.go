package dictionary

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestTerm_StaticCopy(t *testing.T) {
	require := require.New(t)
	testName := "TestTerm_StaticCopy"

	term := Term{}
	err := json.Unmarshal(fixture.Read(t, testName+".json"), &term)
	require.NoError(err)
	test.EmptyFieldsMatch(t, term)
	test.EmptyFieldsMatch(t, term.StaticCopy(), "ID")
}
