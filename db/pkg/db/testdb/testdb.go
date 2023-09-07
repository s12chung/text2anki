// Package testdb contains test helper functions related to db
package testdb

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

// A copy of this constant is in db_test.go
const testDBFile = "testdb.sqlite3"

var tmpPath string

func dbPathF() string    { return path.Join(tmpPath, testDBFile) }
func dbSHAPathF() string { return path.Join(tmpPath, "testdb.sha.txt") }
func init() {
	_, callerFilePath, _, ok := runtime.Caller(0)
	if !ok {
		slog.Error("runtime.Caller not ok for testdb package")
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

// SourcePartMediaImageKey is SourcePartMedia.ImageKey used
const SourcePartMediaImageKey = "testdb.SourcePartMediaImageKey.png"

// Transaction is a wrapper around db.Tx
type Transaction struct{ db.Tx }

// Finalize does nothing (server does not handle Tx, the tests do)
func (t Transaction) Finalize() error { return nil }

// FinalizeError does nothing (server does not handle Tx, the tests do)
func (t Transaction) FinalizeError() error { return nil }

// NewTransaction returns a new Transaction
func NewTransaction(tx db.Tx) Transaction { return Transaction{Tx: tx} }

// TxQs returns a db.NewTxQs used for testing
func TxQs(t *testing.T, opts *sql.TxOptions) db.TxQs {
	require := require.New(t)

	txQs, err := db.NewTxQs(context.Background(), opts)
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
		slog.Error("testdb.MustSetup()", logg.Err(err))
		os.Exit(-1)
	}
}

// Create creates the test database
func Create() error {
	return runWithMatchingDB(func(txQs db.TxQs) error {
		if err := txQs.Create(txQs.Ctx()); err != nil {
			return err
		}
		return SeedList(txQs, nil)
	})
}

func runWithMatchingDB(f func(txQs db.TxQs) error) error {
	dbPath := dbPathF()
	dbSHAPath := dbSHAPathF()

	if err := os.MkdirAll(tmpPath, ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	dbSHA, reuseDB, err := ensureMatchingDB(dbPath, dbSHAPath)
	if err != nil {
		return err
	}
	if !reuseDB {
		if err := db.SetDB(dbPath); err != nil {
			return err
		}
		txQs, err := db.NewTxQs(context.Background(), db.WriteOpts())
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
	return os.WriteFile(dbSHAPath, []byte(dbSHA), ioutil.OwnerRWGroupR)
}

func ensureMatchingDB(dbPath, dbSHAPath string) (string, bool, error) {
	shaSlice := make([]string, len(seederMap)+1)
	shaSlice[0] = fmt.Sprintf("%x", sha256.Sum256(db.SchemaBytes()))
	i := 1
	for _, s := range seederMap {
		bytes, err := s.ReadFile()
		if err != nil {
			return "", false, err
		}
		shaSlice[i] = fmt.Sprintf("%x", sha256.Sum256(bytes))
		i++
	}
	sort.Strings(shaSlice)
	sha := strings.Join(shaSlice, "\n")

	if _, err := os.Stat(dbPath); err != nil {
		//nolint:nilerr // skip if dbPath it doesn't exist
		return sha, false, nil
	}

	//nolint:gosec // for tests, constant path
	shaBytes, err := os.ReadFile(dbSHAPath)
	if err != nil {
		return sha, false, os.Remove(dbPath)
	}
	if string(shaBytes) == sha {
		return sha, true, nil
	}
	return sha, false, os.Remove(dbPath)
}
