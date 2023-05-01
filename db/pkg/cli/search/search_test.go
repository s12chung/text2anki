package search

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/pkg/util/ioutils"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestTermsSearchToCSVRows(t *testing.T) {
	require := require.New(t)
	testName := "TestTermsSearchToCSVRows"
	testdb.SetupTempDBT(t, testName)
	testdb.Seed(t)

	terms, err := db.Qs().TermsSearch(context.Background(), testdb.SearchTerm, testdb.SearchConfig)
	require.NoError(err)
	rows, err := TermsSearchToCSVRows(terms)
	require.NoError(err)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, rows))
}

var testConfig = Config{
	Queries: []string{"a", "b"},
	Config:  testdb.SearchConfig,
}
var changedTestConfig = Config{
	Queries: []string{"test", "change"},
	Config: db.TermsSearchConfig{
		PopLog:       1,
		PopWeight:    2,
		CommonWeight: 3,
		LenLog:       4,
	},
}

func TestConfigToCSVRows(t *testing.T) {
	testName := "TestConfigToCSVRows"
	rows := ConfigToCSVRows(testConfig)
	fixture.CompareReadOrUpdate(t, testName+".json", fixture.JSON(t, rows))
}

func TestGetOrDefaultConfig(t *testing.T) {
	oldConfig := defaultConfig
	defaultConfig = testConfig
	defer func() {
		defaultConfig = oldConfig
	}()

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

	err = os.WriteFile(configPath, fixture.JSON(t, changedTestConfig), ioutils.OwnerRWGroupR)
	require.NoError(err)
	config, err = GetOrDefaultConfig(configPath)
	require.NoError(err)
	require.Equal(string(fixture.JSON(t, changedTestConfig)), string(fixture.JSON(t, config)))
}
