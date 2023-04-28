// Package db provides functions related to the database
package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sql.Open needs it from init()
)

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

//go:embed schema.sql
var schema string

// Create creates the tables from schema.sql
func Create(ctx context.Context) error {
	if _, err := DB().ExecContext(ctx, schema); err != nil {
		return err
	}
	return nil
}

//go:embed custom/TermsSearch.sql
var termsSearch string

// TermsSearchRow is the row returned by TermsSearch
type TermsSearchRow struct {
	Term
	Rank     sql.NullFloat64
	CalcRank sql.NullFloat64
}

// TermsSearch searches within Terms for txt using fts5
func (q *Queries) TermsSearch(ctx context.Context, text string) ([]TermsSearchRow, error) {
	textLength := len(text)
	search := fmt.Sprintf("%v OR %v*", text, text)

	rows, err := q.db.QueryContext(ctx, termsSearch, textLength, textLength, search)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // it's fine, just closing row
	defer rows.Close()
	var items []TermsSearchRow
	for rows.Next() {
		var i TermsSearchRow
		if err := rows.Scan(
			&i.Text,
			&i.Variants,
			&i.PartOfSpeech,
			&i.CommonLevel,
			&i.Translations,
			&i.Popularity,
			&i.Rank,
			&i.CalcRank,
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
