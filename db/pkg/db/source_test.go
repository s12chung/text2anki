package db

import (
	"bytes"
	"context"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizer"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func firstSource(t *testing.T) Source {
	require := require.New(t)
	source, err := Qs().SourceGet(context.Background(), 1)
	require.NoError(err)
	return source
}

func TestSourceStructured_StaticCopy(t *testing.T) {
	require := require.New(t)
	test.EmptyFieldsMatch(t, firstSource(t))

	sourceCopy := firstSource(t).ToSourceStructured()
	sourceCopy.ID = 0
	sourceCopy.UpdatedAt = time.Time{}
	sourceCopy.CreatedAt = time.Time{}
	require.Equal(sourceCopy, firstSource(t).ToSourceStructured().StaticCopy())
}

func TestSourcePartMedia_MarshalJSON(t *testing.T) {
	testName := "TestSourcePartMedia_MarshalJSON"

	testCases := []struct {
		name             string
		prepareSerialize bool
	}{
		{name: "basic"},
		{name: "prepare_serialize", prepareSerialize: true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			source := firstSource(t).ToSourceStructured()
			source.Parts = setupParts(t, source.Parts[0], testUUID)
			if tc.prepareSerialize {
				source.PrepareSerialize()
			}
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, source.StaticCopy()))
		})
	}
}

func setupParts(t *testing.T, part SourcePart, prePartListID string) []SourcePart {
	require := require.New(t)

	parts := make([]SourcePart, 3)
	baseKey := storage.BaseKey(SourcesTable, PartsColumn, prePartListID)
	for i := 0; i < len(parts); i++ {
		parts[i] = part

		key := baseKey + ".PreParts[" + strconv.Itoa(i) + "].Image.txt"
		parts[i].Media = &SourcePartMedia{ImageKey: key}
		require.NoError(storageAPI.Store(key, bytes.NewReader([]byte("image"+strconv.Itoa(i)))))
	}
	for i := 0; i < 1; i++ {
		key := baseKey + ".PreParts[0].Audio.txt"
		parts[i].Media.AudioKey = key
		require.NoError(storageAPI.Store(key, bytes.NewReader([]byte("audio"+strconv.Itoa(i)+"!"))))
	}
	return parts
}

func TestSourcePartMedia_UnmarshalJSON(t *testing.T) {
	t.Skip("Function should be removed if tests are improved such that they're not relying on Unmarshalling requests")
}

func TestSourceStructured_DefaultedName(t *testing.T) {
	require := require.New(t)

	source := firstSource(t).ToSourceStructured()
	require.Equal(source.Name, source.DefaultedName())
	source.Name = ""
	require.Equal(source.Parts[0].TokenizedTexts[0].Text.Text, source.DefaultedName())
}

func TestSourceStructured_UpdateParams(t *testing.T) {
	testName := "TestSourceStructured_UpdateParams"
	test.EmptyFieldsMatch(t, firstSource(t))
	createParams := firstSource(t).ToSourceStructured().UpdateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestSourceStructured_CreateParams(t *testing.T) {
	testName := "TestSourceStructured_CreateParams"
	test.EmptyFieldsMatch(t, firstSource(t))
	createParams := firstSource(t).ToSourceStructured().CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestSource_ToSource_ToSourceStructured(t *testing.T) {
	test.EmptyFieldsMatch(t, firstSource(t))
	test.EmptyFieldsMatch(t, firstSource(t).ToSourceStructured())
	reflect.DeepEqual(firstSource(t), firstSource(t).ToSourceStructured().ToSource())
}

var textTokenizer = TextTokenizer{
	Parser:       text.NewParser(text.Korean, text.English),
	Tokenizer:    tokenizer.NewSplitTokenizer(),
	CleanSpeaker: true,
}

func TestTextTokenizer_TokenizedTexts(t *testing.T) {
	testName := "TestTextTokenizer_TokenizedTexts"

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

			s := string(fixture.Read(t, path.Join(testName, tc.name+".txt")))
			split := strings.Split(s, "===")
			if len(split) == 1 {
				split = append(split, "")
			}
			tokenizedTexts, err := textTokenizer.TokenizedTexts(split[0], split[1])
			require.NoError(err)

			nonSpeaker := strings.TrimPrefix(tc.name, "speaker_")
			fixture.CompareReadOrUpdate(t, path.Join(testName, nonSpeaker+".json"), fixture.JSON(t, tokenizedTexts))
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

	txQs := TxQsT(t)
	source, err := txQs.SourceCreate(txQs.Ctx(), firstSource(t).ToSourceStructured().CreateParams())
	require.NoError(err)
	testRecentTimestamps(t, source.CreatedAt, source.UpdatedAt)
}

func TestQueries_SourceUpdate(t *testing.T) {
	require := require.New(t)

	txQs := TxQsT(t)

	newSource, err := txQs.SourceCreate(txQs.Ctx(), firstSource(t).ToSourceStructured().CreateParams())
	require.NoError(err)
	time.Sleep(1 * time.Second)

	source, err := txQs.SourceUpdate(txQs.Ctx(), newSource.ToSourceStructured().UpdateParams())
	require.NoError(err)
	testRecentTimestamps(t, source.UpdatedAt)
	require.NotEqual(newSource.UpdatedAt, source.UpdatedAt)
}
