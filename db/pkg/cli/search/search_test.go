package search

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func init() { testdb.MustSetup() }

func TestTermsSearchToCSVRows(t *testing.T) {
	require := require.New(t)
	testName := "TestTermsSearchToCSVRows"

	txQs := testdb.TxQs(t, nil)

	terms, err := txQs.TermsSearchRaw(txQs.Ctx(), testdb.SearchTerm, testdb.SearchPOS, testdb.SearchConfig)
	require.NoError(err)
	rows, err := TermsSearchToCSVRows(terms)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, rows))

	_, err = txQs.TermsSearchRaw(txQs.Ctx(), testdb.SearchTerm, testdb.SearchPOS, db.TermsSearchConfig{
		PosWeight:    30,
		PopWeight:    40,
		CommonWeight: 40,
	})
	require.Error(err)
}

var testConfig = Config{
	Queries: []Query{{Str: "a"}, {Str: "b"}},
}
var changedTestConfig = Config{
	Queries: []Query{{Str: "test"}, {Str: "change"}},
}

func TestConfigToCSVRows(t *testing.T) {
	testName := "TestConfigToCSVRows"
	rows := ConfigToCSVRows()
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, rows))
}

func TestGetOrDefaultConfig(t *testing.T) {
	oldConfig := defaultConfig
	defaultConfig = testConfig
	t.Cleanup(func() { defaultConfig = oldConfig })

	require := require.New(t)
	testName := "TestGetOrDefaultConfig"
	configPath := path.Join(os.TempDir(), test.GenerateFilename(testName, ".json"))

	// no file exists
	config, err := GetOrDefaultConfig(configPath)
	require.NoError(err)
	require.Equal(Config{}, config)
	// check the default config
	//nolint:gosec // it's tests
	fileConfig, err := os.ReadFile(configPath)
	require.NoError(err)
	require.Equal(string(fixture.JSON(t, testConfig)), string(fileConfig))

	err = os.WriteFile(configPath, fixture.JSON(t, changedTestConfig), ioutil.OwnerRWGroupR)
	require.NoError(err)
	config, err = GetOrDefaultConfig(configPath)
	require.NoError(err)
	require.Equal(string(fixture.JSON(t, changedTestConfig)), string(fixture.JSON(t, config)))
}
