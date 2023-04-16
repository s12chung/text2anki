// Package db provides functions related to the database
package db

import (
	"database/sql"
)

var database *sql.DB

func init() {
	var err error
	database, err = sql.Open("sqlite3", "data.sqlite3")
	if err != nil {
		panic("database/sql.Open error: " + err.Error())
	}
}

// DB returns the database
func DB() *sql.DB {
	return database
}
