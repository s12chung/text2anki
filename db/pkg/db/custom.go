// Package db provides functions related to the database
package db

import (
	"context"
	"database/sql"
	_ "embed"

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
var schema string

// Create creates the tables from schema.sql
func Create(ctx context.Context) error {
	if _, err := DB().ExecContext(ctx, schema); err != nil {
		return err
	}
	return nil
}
