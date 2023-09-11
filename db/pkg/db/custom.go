// Package db provides functions related to the database
package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3" // sql.Open needs it from init()

	"github.com/s12chung/text2anki/pkg/storage"
)

const arraySeparator = ", "

var database *sql.DB

// SetDB sets the database for the global database
func SetDB(dataSourceName string) error {
	var err error
	// related to require above
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	return nil
}

// Tx represents a transaction that also carries the context
type Tx interface {
	DBTX
	Ctx() context.Context

	Commit() error
	Rollback() error

	Finalize() error
	FinalizeError() error
}

// Transaction represents a transaction that also carries the context and implements reqtx.Tx
type Transaction struct {
	*sql.Tx
	ctx context.Context //nolint:containedctx //it's very clear what the context is about, the transaction
}

// Ctx returns the context of the transaction
func (t Transaction) Ctx() context.Context { return t.ctx }

// Finalize commits the Tx
func (t Transaction) Finalize() error { return t.Commit() }

// FinalizeError rolls back the Tx
func (t Transaction) FinalizeError() error { return t.Rollback() }

// TxQs is a queries with a transaction attached
type TxQs struct {
	*Queries
	Tx
}

// WriteOpts returns write *sql.TxOptions
func WriteOpts() *sql.TxOptions { return &sql.TxOptions{} }

// NewTxQs returns a new TxQs
func NewTxQs(ctx context.Context, opts *sql.TxOptions) (TxQs, error) {
	if opts == nil {
		opts = &sql.TxOptions{ReadOnly: true}
	}
	tx, err := database.BeginTx(ctx, opts)
	if err != nil {
		return TxQs{}, err
	}
	return TxQs{Tx: &Transaction{Tx: tx, ctx: ctx}, Queries: New(tx)}, nil
}

//go:embed schema.sql
var schemaSQL string

//go:embed schema.sql
var schemaBytes []byte

// SchemaBytes returns the bytes of the schema
func SchemaBytes() []byte {
	return schemaBytes
}

// Create creates the tables from schema.sql
func (q *Queries) Create(ctx context.Context) error {
	if _, err := q.db.ExecContext(ctx, schemaSQL); err != nil {
		return err
	}
	return nil
}

//go:embed custom/TableNames.sql
var tableNamesSQL string

// TableNames returns all the table names
func (q *Queries) TableNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, tableNamesSQL)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // it's fine, just closing row
	defer rows.Close()
	var items []string
	for rows.Next() {
		var i string
		if err := rows.Scan(
			&i,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// ClearAll clears all data in the tables
func (q *Queries) ClearAll(ctx context.Context) error {
	tableNames, err := q.TableNames(ctx)
	if err != nil {
		return err
	}

	sql := ""
	for _, tableName := range tableNames {
		sql += fmt.Sprintf("DELETE FROM %v; ", tableName)
	}

	if _, err := q.db.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}

var dbStorage storage.DBStorage

// SetDBStorage sets the storage.DBStorage used in model JSON marshall/unmarshall
func SetDBStorage(d storage.DBStorage) { dbStorage = d }

var plog *slog.Logger

// SetLog setts the log for the package
func SetLog(log *slog.Logger) { plog = log }
