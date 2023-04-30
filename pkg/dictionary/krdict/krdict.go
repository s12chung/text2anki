// Package krdict contains functions for a Korean dictionary connected to a database
package krdict

import (
	"context"
	"database/sql"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/dictionary"
)

// New returns a new KrDict
func New(database *sql.DB) *KrDict {
	return &KrDict{db: database, queries: db.New(database)}
}

// KrDict is a Korean dictionary connected to a database
type KrDict struct {
	db *sql.DB

	queries *db.Queries
}

// Search searches for the query inside the dictionary
func (k *KrDict) Search(q string) ([]dictionary.Term, error) {
	rows, err := k.queries.TermsSearch(context.Background(), q, db.DefaultTermsSearchConfig())
	if err != nil {
		return nil, err
	}
	terms := make([]dictionary.Term, len(rows))
	for i, row := range rows {
		term, err := row.Term.DictionaryTerm()
		if err != nil {
			return nil, err
		}
		terms[i] = term
	}
	return terms, nil
}
