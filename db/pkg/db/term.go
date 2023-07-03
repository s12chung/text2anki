package db

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/stringutil"
)

// ToDBTerm converts a dictionary.Term to a Term
func ToDBTerm(term dictionary.Term, popularity int) (Term, error) {
	variants := strings.Join(term.Variants, arraySeparator)

	translations, err := json.Marshal(term.Translations)
	if err != nil {
		return Term{}, err
	}

	return Term{
		ID:           term.ID,
		Text:         term.Text,
		Variants:     variants,
		PartOfSpeech: string(term.PartOfSpeech),
		CommonLevel:  int64(term.CommonLevel),
		Translations: string(translations),
		Popularity:   int64(popularity),
	}, nil
}

// DictionaryTerm converts the term to a dictionary.Term
func (t Term) DictionaryTerm() (dictionary.Term, error) {
	variants := stringutil.SplitClean(t.Variants, arraySeparator)
	var translations []dictionary.Translation
	if err := json.Unmarshal([]byte(t.Translations), &translations); err != nil {
		return dictionary.Term{}, err
	}

	return dictionary.Term{
		ID:               t.ID,
		Text:             t.Text,
		Variants:         variants,
		PartOfSpeech:     lang.PartOfSpeech(t.PartOfSpeech),
		CommonLevel:      lang.CommonLevel(t.CommonLevel),
		Translations:     translations,
		DictionarySource: "Korean Basic Dictionary (23-03)",
	}, nil
}

// CreateParams converts the term to a TermCreateParams
func (t Term) CreateParams() TermCreateParams {
	return TermCreateParams{
		Text:         t.Text,
		Variants:     t.Variants,
		PartOfSpeech: t.PartOfSpeech,
		CommonLevel:  t.CommonLevel,
		Translations: t.Translations,
		Popularity:   t.Popularity,
	}
}

// TermsSearchRow is the row returned by TermsSearch
type TermsSearchRow struct {
	Term       `json:"term"`
	PosCalc    sql.NullFloat64 `json:"pos_calc"`
	PopCalc    sql.NullFloat64 `json:"pop_calc"`
	CommonCalc sql.NullFloat64 `json:"common_calc"`
	LenCalc    sql.NullFloat64 `json:"len_calc"`
}

// TermsSearchConfig is the config for TermsSearchRaw
type TermsSearchConfig struct {
	PosWeight    int `json:"pos_weight,omitempty"`
	PopLog       int `json:"pop_log,omitempty"`
	PopWeight    int `json:"pop_weight,omitempty"`
	CommonWeight int `json:"common_weight,omitempty"`
	LenLog       int `json:"len_log,omitempty"`
}

var defaultTermsSearchConfig = TermsSearchConfig{
	PosWeight:    15,
	PopLog:       50,
	PopWeight:    25,
	CommonWeight: 10,
	LenLog:       2,
}

// DefaultTermsSearchConfig returns the default TermsSearchConfig
func DefaultTermsSearchConfig() TermsSearchConfig {
	return defaultTermsSearchConfig
}

// TermsSearch searches within Terms for text given the default config
func (q *Queries) TermsSearch(ctx context.Context, query string, pos lang.PartOfSpeech) ([]TermsSearchRow, error) {
	return q.TermsSearchRaw(ctx, query, pos, DefaultTermsSearchConfig())
}

//go:embed custom/TermsSearch.sql
var termsSearchSQL string

// TermsSearchRaw searches within Terms for text given the config
func (q *Queries) TermsSearchRaw(ctx context.Context, query string, pos lang.PartOfSpeech, c TermsSearchConfig) ([]TermsSearchRow, error) {
	if c.PosWeight+c.PopWeight+c.CommonWeight > 100 {
		return nil, fmt.Errorf("c.PosWeight + config.PopWeight + config.CommonWeight > 100: %v, %v, %v", c.PosWeight, c.PopWeight, c.CommonWeight)
	}

	rows, err := q.db.QueryContext(ctx, termsSearchSQL, query, pos, c.PosWeight, c.PopLog, c.PopWeight, c.CommonWeight, c.LenLog)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // it's fine, just closing row
	defer rows.Close()
	var items []TermsSearchRow
	for rows.Next() {
		var i TermsSearchRow
		if err := rows.Scan(
			&i.ID,
			&i.Text,
			&i.Variants,
			&i.PartOfSpeech,
			&i.CommonLevel,
			&i.Translations,
			&i.Popularity,
			&i.PosCalc,
			&i.PopCalc,
			&i.CommonCalc,
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
