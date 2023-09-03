package db

import (
	"encoding/json"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestToDBTerm(t *testing.T) {
	require := require.New(t)
	testName := "TestToDBTerm"

	term := dictionary.Term{}
	err := json.Unmarshal(fixture.Read(t, testName+"Src.json"), &term)
	require.NoError(err)
	test.EmptyFieldsMatch(t, term)

	dbTerm, err := ToDBTerm(term, 1)
	require.NoError(err)
	test.EmptyFieldsMatch(t, dbTerm)

	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, dbTerm))
}

func TestTerm_DictionaryTerm(t *testing.T) {
	testName := "TestTerm_DictionaryTerm"
	tcs := []struct {
		name        string
		emptyFields []string
	}{
		{name: "Base"},
		{name: "EmptyVariants", emptyFields: []string{"Variants"}},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			term := Term{}
			err := json.Unmarshal(fixture.Read(t, path.Join(testName, tc.name+"Src.json")), &term)
			require.NoError(err)
			test.EmptyFieldsMatch(t, term, tc.emptyFields...)

			dictTerm, err := term.DictionaryTerm()
			require.NoError(err)
			test.EmptyFieldsMatch(t, dictTerm)

			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, dictTerm))
		})
	}
}

func TestTerm_CreateParams(t *testing.T) {
	require := require.New(t)
	testName := "TestTerm_CreateParams"

	term := Term{}
	err := json.Unmarshal(fixture.Read(t, "TestToDBTerm.json"), &term)
	require.NoError(err)
	test.EmptyFieldsMatch(t, term)

	createParams := term.CreateParams()
	test.EmptyFieldsMatch(t, createParams)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, createParams))
}

func TestQueries_TermsSearchRaw(t *testing.T) {
	require := require.New(t)
	testName := "TestQueries_TermsSearchRaw"

	txQs := TxQsT(t)
	searchTerm := "마음"
	searchPOS := lang.PartOfSpeechVerb
	searchConfig := TermsSearchConfig{
		PosWeight:    10,
		PopLog:       20,
		PopWeight:    40,
		CommonWeight: 40,
		LenLog:       2,
	}
	results, err := txQs.TermsSearchRaw(txQs.Ctx(), searchTerm, lang.PartOfSpeechUnknown, searchConfig)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, path.Join(testName, "unknown.json"), fixture.JSON(t, results))

	_, err = txQs.TermsSearchRaw(txQs.Ctx(), searchTerm, lang.PartOfSpeechUnknown, TermsSearchConfig{
		PopWeight:    50,
		CommonWeight: 51,
	})
	require.Error(err)

	results, err = txQs.TermsSearchRaw(txQs.Ctx(), searchTerm, searchPOS, searchConfig)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, path.Join(testName, "verb.json"), fixture.JSON(t, results))

	configCopy := searchConfig
	configCopy.Limit = 1
	results, err = txQs.TermsSearchRaw(txQs.Ctx(), searchTerm, searchPOS, configCopy)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, path.Join(testName, "limit.json"), fixture.JSON(t, results))
}
