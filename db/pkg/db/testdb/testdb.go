// Package testdb contains test helper functions related to db
package testdb

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/models"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

var tmpPath string

func dbPathF() string    { return path.Join(tmpPath, "testdb.sqlite3") }
func dbSHAPathF() string { return path.Join(tmpPath, "testdb.sha.txt") }
func init() {
	_, callerFilePath, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("runtime.Caller not ok for Seed()")
		os.Exit(-1)
	}
	callerPath := path.Dir(callerFilePath)
	tmpPath = path.Join(callerPath, "..", "..", "..", "tmp")
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

// Transaction is a wrapper around db.Tx
type Transaction struct{ db.Tx }

// Finalize does nothing (server does not handle Tx, the tests do)
func (t Transaction) Finalize() error { return nil }

// FinalizeError does nothing (server does not handle Tx, the tests do)
func (t Transaction) FinalizeError() error { return nil }

// NewTransaction returns a new Transaction
func NewTransaction(tx db.Tx) Transaction { return Transaction{Tx: tx} }

// TxQs returns a db.NewTxQs used for testing
func TxQs(t *testing.T) db.TxQs {
	require := require.New(t)

	txQs, err := db.NewTxQs()
	require.NoError(err)

	txQs.Tx = NewTransaction(txQs.Tx)
	t.Cleanup(func() {
		require.NoError(txQs.Rollback())
	})
	return txQs
}

// MustSetup sets up the test database
func MustSetup() {
	if err := db.SetDB(dbPathF()); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Create creates the test database
func Create() error {
	return runWithSafeSchema(func(txQs db.TxQs) error {
		if err := txQs.Create(txQs.Ctx()); err != nil {
			return err
		}
		return models.SeedList(txQs, nil)
	})
}

func runWithSafeSchema(f func(txQs db.TxQs) error) error {
	dbPath := dbPathF()
	dbSHAPath := dbSHAPathF()

	if err := os.MkdirAll(tmpPath, ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	schemaSHA, reuseSchema, err := ensureSafeSchema(dbPath, dbSHAPath)
	if err != nil {
		return err
	}
	if !reuseSchema {
		if err := db.SetDB(dbPath); err != nil {
			return err
		}
		txQs, err := db.NewTxQs()
		if err != nil {
			return err
		}
		defer txQs.Rollback() //nolint:errcheck // rollback can fail if committed
		if err := f(txQs); err != nil {
			return err
		}
		if err := txQs.Commit(); err != nil {
			return err
		}
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
