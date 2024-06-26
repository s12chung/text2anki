package db

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
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
	txQs := TxQsT(t, nil)
	source := firstSource(t, txQs)
	test.EmptyFieldsMatch(t, source)
	test.EmptyFieldsMatch(t, source.ToSourceStructured().StaticCopy(), "ID", "UpdatedAt", "CreatedAt")
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
		t.Run(tc.name, func(t *testing.T) {
			txQs := TxQsT(t, nil)

			source := firstSource(t, txQs).ToSourceStructured()
			source.Parts = setupParts(t, source.Parts[0], testUUID)
			if tc.prepareSerialize {
				source.PrepareSerialize()
			}
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, tc.name), source.StaticCopy())
		})
	}
}

func setupParts(t *testing.T, part SourcePart, prePartListID string) []SourcePart {
	require := require.New(t)

	parts := make([]SourcePart, 3)
	baseKey := storage.BaseKey(SourcesTable, PartsColumn, prePartListID)
	for i := range len(parts) {
		parts[i] = part

		key := baseKey + ".PreParts[" + strconv.Itoa(i) + "].Image.txt"
		parts[i].Media = &SourcePartMedia{ImageKey: key}
		require.NoError(storageAPI.Store(key, bytes.NewReader([]byte("image"+strconv.Itoa(i)))))
	}
	for i := range 1 {
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

	source.Parts[0].TokenizedTexts = nil
	require.Equal("", source.DefaultedName())

	source.Parts = nil
	require.Equal("", source.DefaultedName())
}

func TestSourceStructured_UpdateParams(t *testing.T) {
	testName := "TestSourceStructured_UpdateParams"
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	updateParams := firstSource(t, txQs).ToSourceStructured().UpdateParams()
	test.EmptyFieldsMatch(t, updateParams)
	fixture.CompareReadOrUpdateJSON(t, testName, updateParams)
}

func TestSourceStructured_UpdatePartsParams(t *testing.T) {
	testName := "TestSourceStructured_UpdatePartsParams"
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	updateParams := firstSource(t, txQs).ToSourceStructured().UpdatePartsParams()
	test.EmptyFieldsMatch(t, updateParams)
	fixture.CompareReadOrUpdateJSON(t, testName, updateParams)
}

func TestSourceStructured_CreateParams(t *testing.T) {
	testName := "TestSourceStructured_CreateParams"
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	createParams := firstSource(t, txQs).ToSourceStructured().CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdateJSON(t, testName, createParams)
}

func TestSource_ToSource_ToSourceStructured(t *testing.T) {
	txQs := TxQsT(t, nil)

	test.EmptyFieldsMatch(t, firstSource(t, txQs))
	test.EmptyFieldsMatch(t, firstSource(t, txQs).ToSourceStructured())
	reflect.DeepEqual(firstSource(t, txQs), firstSource(t, txQs).ToSourceStructured().ToSource())
}

type crcTranslator struct{}

func (c crcTranslator) Translate(_ context.Context, s string) (string, error) {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("crc-%x", crc32.ChecksumIEEE([]byte(line)))
	}
	return strings.Join(lines, "\n"), nil
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
		{name: "translator"},
		{name: "translator_weave"},
		{name: "translator_none"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			t.Parallel()

			textTokenizerDup := textTokenizer
			if strings.HasPrefix(tc.name, "translator") {
				textTokenizerDup.Translator = crcTranslator{}
			}

			split := strings.Split(string(fixture.Read(t, path.Join(testName, tc.name+".txt"))), "===")
			if len(split) == 1 {
				split = append(split, "")
			}

			tokenizedTexts, err := textTokenizerDup.TokenizedTexts(context.Background(), split[0], split[1])
			require.NoError(err)

			nonSpeaker := strings.TrimPrefix(tc.name, "speaker_")
			mutex.Lock()
			t.Cleanup(mutex.Unlock)
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, nonSpeaker), tokenizedTexts)
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

func TestQueries_SourceStructuredIndex(t *testing.T) {
	testName := "TestQueries_SourceStructuredIndex"

	testCases := []struct {
		name string
	}{
		{name: "basic"},
		{name: "clear"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			txQs := TxQsT(t, nil)
			if tc.name == "clear" {
				txQs = TxQsT(t, WriteOpts())
				require.NoError(txQs.ClearAllTable(txQs.Ctx(), "sources"))
			}

			sourceStructureds, err := txQs.SourceStructuredIndex(txQs.Ctx())
			require.NoError(err)

			var staticCopy []SourceStructured
			if len(sourceStructureds) > 0 {
				staticCopy = make([]SourceStructured, len(sourceStructureds))
				for i := range sourceStructureds {
					staticCopy[i] = sourceStructureds[i].StaticCopy()
				}
			}
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, tc.name), staticCopy)
		})
	}
}
