// Package db provides functions related to the database
package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sql.Open needs it from init()
)

const arraySeparator = ", "

var database *sql.DB

// SetDB sets the database returned from the DB() function
func SetDB(dataSourceName string) error {
	var err error
	// related to require above
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	return nil
}

// DB returns the database set by SetDB()
func DB() *sql.DB {
	return database
}

// Qs returns the Queries for the database returned from the DB() function
func Qs() *Queries {
	return &Queries{db: DB()}
}

//go:embed schema.sql
var schemaSQL string

// Create creates the tables from schema.sql
func (q *Queries) Create(ctx context.Context) error {
	if _, err := DB().ExecContext(ctx, schemaSQL); err != nil {
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

	if _, err := DB().ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}
