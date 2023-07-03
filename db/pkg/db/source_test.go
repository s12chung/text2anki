package db_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizers"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func firstSource(t *testing.T) db.Source {
	require := require.New(t)
	source, err := db.Qs().SourceGet(context.Background(), 1)
	require.NoError(err)
	return source
}

func TestSourceSerialized_StaticCopy(t *testing.T) {
	testName := "TestSourceSerialized_StaticCopy"
	test.EmptyFieldsMatch(t, firstSource(t))

	staticCopy := firstSource(t).ToSourceSerialized().StaticCopy()
	test.EmptyFieldsMatch(t, staticCopy, "ID", "UpdatedAt", "CreatedAt")
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, staticCopy))
}

func TestSourceSerialized_UpdateParams(t *testing.T) {
	testName := "TestSourceSerialized_UpdateParams"
	test.EmptyFieldsMatch(t, firstSource(t))
	createParams := firstSource(t).ToSourceSerialized().UpdateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestSourceSerialized_CreateParams(t *testing.T) {
	testName := "TestSourceSerialized_CreateParams"
	test.EmptyFieldsMatch(t, firstSource(t))
	createParams := firstSource(t).ToSourceSerialized().CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestSource_ToSource_ToSourceSerialized(t *testing.T) {
	test.EmptyFieldsMatch(t, firstSource(t))
	test.EmptyFieldsMatch(t, firstSource(t).ToSourceSerialized())
	reflect.DeepEqual(firstSource(t), firstSource(t).ToSourceSerialized().ToSource())
}

var textTokenizer = db.TextTokenizer{
	Parser:       text.NewParser(text.Korean, text.English),
	Tokenizer:    tokenizers.NewSplitTokenizer(),
	CleanSpeaker: true,
}

func TestTextTokenizer_TokenizedTexts(t *testing.T) {
	testNamePath := "TestTextTokenizer_TokenizedTexts/"

	testCases := []struct {
		name string
	}{
		{name: "split"},
		{name: "weave"},
		{name: "speaker_split"},
		{name: "speaker_weave"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			s := string(fixture.Read(t, testNamePath+tc.name+".txt"))
			split := strings.Split(s, "===")
			if len(split) == 1 {
				split = append(split, "")
			}
			tokenizedTexts, err := textTokenizer.TokenizedTexts(split[0], split[1])
			require.NoError(err)

			nonSpeaker := strings.TrimPrefix(tc.name, "speaker_")
			fixture.CompareReadOrUpdate(t, testNamePath+nonSpeaker+".json", fixture.JSON(t, tokenizedTexts))
		})
	}
}

func TestTextTokenizer_TokenizeTexts(t *testing.T) {
	require := require.New(t)

	texts := []text.Text{
		{Text: "내가 가는 이길이", Translation: "The road that I’m taking"},
		{Text: "어디로 가는지", Translation: "Where it’s leading me to, where it’s taking me"},
	}

	tokenizedTexts, err := textTokenizer.TokenizeTexts(texts)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestTextTokenizer_TokenizeTexts.json", fixture.JSON(t, tokenizedTexts))
}

func TestQueries_SourceCreate(t *testing.T) {
	require := require.New(t)
	source, err := db.Qs().SourceCreate(context.Background(), firstSource(t).ToSourceSerialized().CreateParams())
	require.NoError(err)
	testRecentTimestamps(t, source.CreatedAt, source.UpdatedAt)
}

func TestQueries_SourceUpdate(t *testing.T) {
	require := require.New(t)
	t.Parallel()

	newSource, err := db.Qs().SourceCreate(context.Background(), firstSource(t).ToSourceSerialized().CreateParams())
	require.NoError(err)
	time.Sleep(1 * time.Second)

	source, err := db.Qs().SourceUpdate(context.Background(), newSource.ToSourceSerialized().UpdateParams())
	require.NoError(err)
	testRecentTimestamps(t, source.UpdatedAt)
	require.NotEqual(newSource.UpdatedAt, source.UpdatedAt)
}
