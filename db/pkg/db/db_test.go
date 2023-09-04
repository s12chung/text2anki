package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/storage"
	"github.com/s12chung/text2anki/pkg/storage/localstore"
	"github.com/s12chung/text2anki/pkg/util/logg"
	"github.com/s12chung/text2anki/pkg/util/test"
	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

func TestMain(m *testing.M) {
	if err := run(m); err != nil {
		slog.Error("db_test.TestMain", logg.Err(err))
		os.Exit(-1)
	}
}

// A copy of this constant is in testdb.go
const testDBFile = "testdb.sqlite3"
const testUUID = "a1234567-3456-9abc-d123-456789abcdef"

var storageAPI localstore.API

func run(m *testing.M) error {
	if err := SetDB(path.Join("..", "..", "tmp", testDBFile)); err != nil {
		return err
	}
	var err error
	storageAPI, err = newStorageAPI("db_test")
	if err != nil {
		return err
	}
	SetDBStorage(storage.NewDBStorage(storageAPI, nil))

	if err := textTokenizer.Setup(); err != nil {
		return err
	}
	code := m.Run()
	if err := textTokenizer.Cleanup(); err != nil {
		return err
	}
	os.Exit(code)
	return nil
}

func newStorageAPI(dirPrefix string) (localstore.API, error) {
	encryptor, err := localstore.NewAESEncryptorFromFile(fixture.JoinTestData("localstore.key"))
	if err != nil {
		return localstore.API{}, err
	}
	return localstore.NewAPI("http://localhost:3000", path.Join(os.TempDir(), test.GenerateName(dirPrefix)), encryptor), nil
}

func testRecentTimestamps(t *testing.T, timestamps ...time.Time) {
	require := require.New(t)
	for _, timestamp := range timestamps {
		require.Greater(time.Now(), timestamp)
	}
}

type TestTransaction struct{ Tx }

func (t TestTransaction) Finalize() error      { return nil }
func (t TestTransaction) FinalizeError() error { return nil }
func NewTestTransaction(tx Tx) TestTransaction { return TestTransaction{Tx: tx} }

func TxQsT(t *testing.T, opts *sql.TxOptions) TxQs {
	require := require.New(t)

	txQs, err := NewTxQs(context.Background(), opts)
	require.NoError(err)

	txQs.Tx = NewTestTransaction(txQs.Tx)
	t.Cleanup(func() {
		require.NoError(txQs.Rollback())
	})
	return txQs
}
