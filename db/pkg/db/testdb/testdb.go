// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

var callerPath string
var dbPath string
var dbSHAPath string

func init() {
	_, callerFilePath, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("runtime.Caller not ok for Seed()")
		os.Exit(-1)
	}
	callerPath = path.Dir(callerFilePath)
	dbPath = path.Join(callerPath, "..", "..", "..", "tmp", "testdb.sqlite3")
	dbSHAPath = path.Join(callerPath, "..", "..", "..", "tmp", "testdb-schema-sha.txt")
}

// MustSetupAndSeed calls Setup() and Seed(), if it fails, it exits
func MustSetupAndSeed() {
	if err := Setup(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := Seed(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// SetupAndSeedT calls Setup() and Seed()
func SetupAndSeedT(t *testing.T) {
	require := require.New(t)
	SetupT(t)
	require.NoError(Seed())
}

// SetupT setups up an empty db and checks errors
func SetupT(t *testing.T) {
	require := require.New(t)
	require.NoError(Setup())
}

// Setup setups up an empty db
func Setup() error {
	if err := os.MkdirAll(path.Dir(dbPath), ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}

	schemaSHA := fmt.Sprintf("%x", sha256.Sum256(db.SchemaBytes()))
	if _, err := os.Stat(dbPath); err == nil {
		//nolint:gosec // for tests, constant path
		shaBytes, err := os.ReadFile(dbSHAPath)
		if err != nil {
			return err
		}
		if string(shaBytes) != schemaSHA {
			if err = os.Remove(dbPath); err != nil {
				return err
			}
		}
	}

	if err := db.SetDB(dbPath); err != nil {
		return err
	}
	if err := db.Qs().Create(context.Background()); err != nil {
		return err
	}
	if err := db.Qs().ClearAll(context.Background()); err != nil {
		return err
	}

	return os.WriteFile(dbSHAPath, []byte(schemaSHA), ioutil.OwnerRWGroupR)
}

// Seed seeds the database with a small amount of data
func Seed() error {
	return SeedModels()
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
