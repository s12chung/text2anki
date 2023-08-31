// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/models"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

var tmpPath string

func init() {
	_, callerFilePath, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("runtime.Caller not ok for Seed()")
		os.Exit(-1)
	}
	callerPath := path.Dir(callerFilePath)
	tmpPath = path.Join(callerPath, "..", "..", "..", "tmp")
}

// NewTx returns a new transaction
func NewTx(t *testing.T) db.Tx {
	require := require.New(t)

	tx, err := db.NewTx()
	require.NoError(err)
	t.Cleanup(func() {
		require.NoError(tx.Rollback())
	})
	return tx
}

// TxQs returns a db.NewTxQs used for testing
func TxQs(t *testing.T) db.TxQs {
	require := require.New(t)

	txQs, err := db.NewTxQs()
	require.NoError(err)
	t.Cleanup(func() {
		require.NoError(txQs.Rollback())
	})
	return txQs
}

// MustSetupAndSeed calls Setup() and Seed(), if it fails, it exits
func MustSetupAndSeed(packageStruct any) {
	MustSetupPackage(packageStruct)
	if err := Seed(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// MustSetupPackage calls Setup(), if it fails, it exits
func MustSetupPackage(packageStruct any) {
	if err := Setup(path.Base(reflect.TypeOf(packageStruct).PkgPath())); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// MustSetup sets up the main database
func MustSetup() {
	if err := Setup("main"); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := Seed(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// MustSetupEmpty sets up the empty database
func MustSetupEmpty() {
	if err := Setup("empty"); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// SetupAndSeedT calls Setup() and Seed()
func SetupAndSeedT(t *testing.T, testName string) {
	require := require.New(t)
	SetupT(t, testName)
	require.NoError(Seed())
}

// SetupT setups up an empty db and checks errors
func SetupT(t *testing.T, testName string) {
	require := require.New(t)
	require.NoError(Setup(testName))
}

// Setup setups up an empty db
func Setup(name string) error {
	dbPath := path.Join(tmpPath, "testdb."+name+".sqlite3")
	dbSHAPath := path.Join(tmpPath, "testdbsha."+name+".txt")

	if err := os.MkdirAll(path.Dir(dbPath), ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}

	schemaSHA, reuseSchema, err := ensureSafeSchema(dbPath, dbSHAPath)
	if err != nil {
		return err
	}

	if err := db.SetDB(dbPath); err != nil {
		return err
	}
	if !reuseSchema {
		if err := db.Qs().Create(context.Background()); err != nil {
			return err
		}
	}
	if err := db.Qs().ClearAll(context.Background()); err != nil {
		return err
	}

	return os.WriteFile(dbSHAPath, []byte(schemaSHA), ioutil.OwnerRWGroupR)
}

func ensureSafeSchema(dbPath, dbSHAPath string) (string, bool, error) {
	schemaSHA := fmt.Sprintf("%x", sha256.Sum256(db.SchemaBytes()))

	if _, err := os.Stat(dbPath); err != nil {
		//nolint:nilerr // skip if dbPath it doesn't exist
		return schemaSHA, false, nil
	}

	//nolint:gosec // for tests, constant path
	shaBytes, err := os.ReadFile(dbSHAPath)
	if err != nil {
		return schemaSHA, false, err
	}
	if string(shaBytes) == schemaSHA {
		return schemaSHA, true, nil
	}
	return schemaSHA, false, os.Remove(dbPath)
}

// Seed seeds the database with a small amount of data
func Seed() error {
	tx, err := db.NewTx()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // rollback can fail if committed
	if err := models.SeedList(tx, nil); err != nil {
		return err
	}
	return tx.Commit()
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
