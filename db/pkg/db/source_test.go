package db

import (
	"bytes"
	"context"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/text"
	"github.com/s12chung/text2anki/pkg/tokenizer"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func firstSource(t *testing.T, txQs TxQs) Source {
	require := require.New(t)
	source, err := txQs.SourceGet(txQs.Ctx(), 1)
	require.NoError(err)
	return source
}

func TestSourceStructured_StaticCopy(t *testing.T) {
	require := require.New(t)
	txQs := TxQsT(t, nil)
	test.EmptyFieldsMatch(t, firstSource(t, txQs))

	sourceCopy := firstSource(t, txQs).ToSourceStructured()
	sourceCopy.ID = 0
	sourceCopy.UpdatedAt = time.Time{}
	sourceCopy.CreatedAt = time.Time{}
	require.Equal(sourceCopy, firstSource(t, txQs).ToSourceStructured().StaticCopy())
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
			txQs := TxQsT(t, nil)

			source := firstSource(t, txQs).ToSourceStructured()
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
	txQs := TxQsT(t, nil)

	source := firstSource(t, txQs).ToSourceStructured()
	require.Equal(source.Name, source.DefaultedName())
	source.Name = ""
	require.Equal(source.Parts[0].TokenizedTexts[0].Text.Text, source.DefaultedName())
}

func TestSourceStructured_UpdateParams(t *testing.T) {
	testName := "TestSourceStructured_UpdateParams"
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	createParams := firstSource(t, txQs).ToSourceStructured().UpdateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestSourceStructured_CreateParams(t *testing.T) {
	testName := "TestSourceStructured_CreateParams"
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	createParams := firstSource(t, txQs).ToSourceStructured().CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestSource_ToSource_ToSourceStructured(t *testing.T) {
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	test.EmptyFieldsMatch(t, firstSource(t, txQs).ToSourceStructured())
	reflect.DeepEqual(firstSource(t, txQs), firstSource(t, txQs).ToSourceStructured().ToSource())
}

var textTokenizer = TextTokenizer{
	Parser:       text.NewParser(text.Korean, text.English),
	Tokenizer:    tokenizer.NewSplitTokenizer(),
	CleanSpeaker: true,
}

func TestTextTokenizer_TokenizedTexts(t *testing.T) {
	testName := "TestTextTokenizer_TokenizedTexts"
	t.Parallel()
	mutex := &sync.Mutex{}

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
			t.Parallel()

			s := string(fixture.Read(t, path.Join(testName, tc.name+".txt")))
			split := strings.Split(s, "===")
			if len(split) == 1 {
				split = append(split, "")
			}
			tokenizedTexts, err := textTokenizer.TokenizedTexts(context.Background(), split[0], split[1])
			require.NoError(err)

			nonSpeaker := strings.TrimPrefix(tc.name, "speaker_")
			mutex.Lock()
			fixture.CompareReadOrUpdate(t, path.Join(testName, nonSpeaker+".json"), fixture.JSON(t, tokenizedTexts))
			mutex.Unlock()
		})
	}
}

func TestTextTokenizer_TokenizeTexts(t *testing.T) {
	require := require.New(t)
	t.Parallel()

	texts := []text.Text{
		{Text: "내가 가는 이길이", Translation: "The road that I’m taking"},
		{Text: "어디로 가는지", Translation: "Where it’s leading me to, where it’s taking me"},
	}

	tokenizedTexts, err := textTokenizer.TokenizeTexts(context.Background(), texts)
	require.NoError(err)

	fixture.CompareReadOrUpdate(t, "TestTextTokenizer_TokenizeTexts.json", fixture.JSON(t, tokenizedTexts))
}

func TestQueries_SourceCreate(t *testing.T) {
	require := require.New(t)
	txQs := TxQsT(t, WriteOpts())

	source, err := txQs.SourceCreate(txQs.Ctx(), firstSource(t, txQs).ToSourceStructured().CreateParams())
	require.NoError(err)
	testRecentTimestamps(t, source.CreatedAt, source.UpdatedAt)
}

func TestQueries_SourceUpdate(t *testing.T) {
	require := require.New(t)
	t.Parallel()
	txQs := TxQsT(t, WriteOpts())

	newSource, err := txQs.SourceCreate(txQs.Ctx(), firstSource(t, txQs).ToSourceStructured().CreateParams())
	require.NoError(err)
	time.Sleep(time.Second)

	source, err := txQs.SourceUpdate(txQs.Ctx(), newSource.ToSourceStructured().UpdateParams())
	require.NoError(err)
	testRecentTimestamps(t, source.UpdatedAt)
	require.NotEqual(newSource.UpdatedAt, source.UpdatedAt)
}
