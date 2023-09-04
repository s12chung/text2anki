package khaiii

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testPort = 9001

func TestKhaiii_Tokenize(t *testing.T) {
	test.CISkip(t, "can't run C environment in CI")
	t.Parallel()
	require := require.New(t)

	tokenizer := newKhaiii(testPort)
	require.NoError(tokenizer.Setup())
	defer func() { require.NoError(tokenizer.CleanupAndWait()) }()

	tokens, err := tokenizer.Tokenize("대한민국은 민주공화국이다.")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestKhaiii_Tokenize.json", fixture.JSON(t, tokens))
}
