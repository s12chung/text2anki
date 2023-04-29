package db

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
)

// ToDBTerm converts a dictionary.Term to a Term
func ToDBTerm(term dictionary.Term, popularity int) (Term, error) {
	variants := strings.Join(term.Variants, arraySeparator)

	translations, err := json.Marshal(term.Translations)
	if err != nil {
		return Term{}, err
	}

	return Term{
		Text:         term.Text,
		Variants:     variants,
		PartOfSpeech: string(term.PartOfSpeech),
		CommonLevel:  int64(term.CommonLevel),
		Translations: string(translations),
		Popularity:   int64(popularity),
	}, nil
}

// DictionaryTerm converts the term to a dictionary.Term
func (t *Term) DictionaryTerm() (dictionary.Term, error) {
	variants := strings.Split(t.Variants, arraySeparator)

	var translations []dictionary.Translation
	if err := json.Unmarshal([]byte(t.Translations), &translations); err != nil {
		return dictionary.Term{}, err
	}

	return dictionary.Term{
		Text:             t.Text,
		Variants:         variants,
		PartOfSpeech:     lang.PartOfSpeech(t.PartOfSpeech),
		CommonLevel:      lang.CommonLevel(t.CommonLevel),
		Translations:     translations,
		DictionarySource: "Korean Basic Dictionary (23-03)",
	}, nil
}

// CreateParams converts the term to a TermCreateParams
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
	PopCalc sql.NullFloat64
	LenCalc sql.NullFloat64
}

// TermsSearchConfig is the config for TermsSearchRaw
type TermsSearchConfig struct {
	PopLog    int
	PopWeight int
	LenLog    int
}

var defaultTermsSearchConfig = TermsSearchConfig{
	PopLog:    100,
	PopWeight: 40,
	LenLog:    3,
}

// TermsSearch searches within Terms for text
func (q *Queries) TermsSearch(ctx context.Context, query string) ([]TermsSearchRow, error) {
	return q.TermsSearchRaw(ctx, query, defaultTermsSearchConfig)
}

// TermsSearchRaw searches within Terms for text given the config
func (q *Queries) TermsSearchRaw(ctx context.Context, query string, c TermsSearchConfig) ([]TermsSearchRow, error) {
	rows, err := q.db.QueryContext(ctx, termsSearch, query, c.PopLog, c.PopWeight, c.LenLog)
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
			&i.PopCalc,
			&i.LenCalc,
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
