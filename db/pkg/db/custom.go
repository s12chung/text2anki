// Package db provides functions related to the database
package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // sql.Open needs it from init()
)

var defaultDB *sql.DB

func init() {
	var err error
	// related to require above
	defaultDB, err = sql.Open("sqlite3", "data.sqlite3")
	if err != nil {
		panic("database/sql.Open error: " + err.Error())
	}
}

// DefaultDB returns the default database
func DefaultDB() *sql.DB {
	return defaultDB
}
