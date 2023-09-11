package komoran

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/logg"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

const testPort = 9002

var plog = logg.Default()

func TestKomoran_Tokenize(t *testing.T) {
	test.CISkip(t, "can't run java environment in CI")
	t.Parallel()
	require := require.New(t)
	ctx := context.Background()

	tokenizer := newKomoran(ctx, testPort, plog)
	require.NoError(tokenizer.Setup(ctx))
	defer func() { require.NoError(tokenizer.CleanupAndWait()) }()

	tokens, err := tokenizer.Tokenize(ctx, "대한민국은 민주공화국이다.")
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestKomoran_Tokenize.json", fixture.JSON(t, tokens))
}
