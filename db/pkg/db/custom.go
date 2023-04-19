// Package db provides functions related to the database
package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // sql.Open needs it from init()
)

var database *sql.DB

func init() {
	var err error
	// related to require above
	database, err = sql.Open("sqlite3", "data.sqlite3")
	if err != nil {
		panic("database/sql.Open error: " + err.Error())
	}
}

// DB returns the database
func DB() *sql.DB {
	return database
}
