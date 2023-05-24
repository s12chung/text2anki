// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

// MustSetupAndSeed calls Setup() and Seed(), if it fails, it exits
func MustSetupAndSeed(testName string) {
	if err := SetupTempDB(testName); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := Seed(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// SetupAndSeed calls Setup() and Seed()
func SetupAndSeed(t *testing.T, testName string) {
	SetupTempDBT(t, testName)
	SeedT(t)
}

// SetupTempDB calls db.SetDB with a temp file
func SetupTempDB(testName string) error {
	filename := test.GenerateFilename(testName, ".sqlite3")
	if err := db.SetDB(path.Join(os.TempDir(), filename)); err != nil {
		return err
	}
	return db.Create(context.Background())
}

// SetupTempDBT calls SetupTempDB and checks errors
func SetupTempDBT(t *testing.T, testName string) {
	require := require.New(t)
	err := SetupTempDB(testName)
	require.NoError(err)
}

// Seed seeds the database with a small amount of data
func Seed() error {
	_, callerPath, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("runtime.Caller not ok for Seed()")
	}

	queries := db.Qs()

	var terms []db.Term
	if err := unmarshall(callerPath, "TermsSeed", &terms); err != nil {
		return err
	}
	for _, term := range terms {
		if _, err := queries.TermCreate(context.Background(), term.CreateParams()); err != nil {
			return err
		}
	}

	var sourceSerializeds []db.SourceSerialized
	if err := unmarshall(callerPath, "SourcesSeed", &sourceSerializeds); err != nil {
		return err
	}
	for _, sourceSerialized := range sourceSerializeds {
		if _, err := queries.SourceSerializedCreate(context.Background(), sourceSerialized.TokenizedTexts); err != nil {
			return err
		}
	}
	return nil
}

// SeedT seeds the database with a small amount of data
func SeedT(t *testing.T) {
	require := require.New(t)
	require.NoError(Seed())
}

func unmarshall(callerPath, filename string, data any) error {
	bytes, err := os.ReadFile(path.Join(path.Dir(callerPath), fixture.TestDataDir, filename) + ".json")
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, data)
}

// SearchTerm is a search term used for tests
const SearchTerm = "마음"

// SearchPOS is the search POS for tests
const SearchPOS = lang.PartOfSpeechVerb

// SearchConfig is the config used for test searching (so it stays constant)
var SearchConfig = db.TermsSearchConfig{
	PosWeight:    10,
	PopLog:       20,
	PopWeight:    40,
	CommonWeight: 40,
	LenLog:       2,
}
