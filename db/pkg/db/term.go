package db

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/s12chung/text2anki/pkg/dictionary"
)

// ToDBTerm converts a dictionary.Term to a Term
func ToDBTerm(term dictionary.Term, popularity int) (Term, error) {
	variants, err := json.Marshal(term.Variants)
	if err != nil {
		return Term{}, err
	}
	translations, err := json.Marshal(term.Translations)
	if err != nil {
		return Term{}, err
	}

	return Term{
		Text:         term.Text,
		Variants:     string(variants),
		PartOfSpeech: string(term.PartOfSpeech),
		CommonLevel:  strconv.Itoa(int(term.CommonLevel)),
		Translations: string(translations),
		Popularity:   strconv.Itoa(popularity),
	}, nil
}

// CreateParams converts a term to a TermCreateParams
func (t *Term) CreateParams() TermCreateParams {
	return TermCreateParams{
		Text:         t.Text,
		Variants:     t.Variants,
		PartOfSpeech: t.PartOfSpeech,
		CommonLevel:  t.CommonLevel,
		Translations: t.Translations,
		Popularity:   t.Popularity,
	}
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
