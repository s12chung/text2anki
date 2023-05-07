package db

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/s12chung/text2anki/pkg/dictionary"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/stringclean"
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
	variants := stringclean.Split(t.Variants, arraySeparator)
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
	PosCalc    sql.NullFloat64
	PopCalc    sql.NullFloat64
	CommonCalc sql.NullFloat64
	LenCalc    sql.NullFloat64
}

// TermsSearchConfig is the config for TermsSearchRaw
type TermsSearchConfig struct {
	PosWeight    int `json:"pos_weight"`
	PopLog       int `json:"pop_log"`
	PopWeight    int `json:"pop_weight"`
	CommonWeight int `json:"common_weight"`
	LenLog       int `json:"len_log"`
}

func init() {
	firm.RegisterType(firm.NewTypedDefinition(TermsSearchConfig{}).Validates(firm.RuleMap{
		"PopLog":       {rule.Presence{}},
		"PopWeight":    {rule.Presence{}},
		"CommonWeight": {rule.Presence{}},
		"LenLog":       {rule.Presence{}},
	}))
}

var defaultTermsSearchConfig = TermsSearchConfig{
	PosWeight:    50,
	PopLog:       50,
	PopWeight:    30,
	CommonWeight: 15,
	LenLog:       2,
}

// DefaultTermsSearchConfig returns the default TermsSearchConfig
func DefaultTermsSearchConfig() TermsSearchConfig {
	return defaultTermsSearchConfig
}

// TermsSearch searches within Terms for text given the config
func (q *Queries) TermsSearch(ctx context.Context, query string, pos lang.PartOfSpeech, c TermsSearchConfig) ([]TermsSearchRow, error) {
	if c.PopWeight+c.CommonWeight > 100 {
		return nil, fmt.Errorf("config.PopWeight and config.CommonWeight > 100: %v, %v", c.PopWeight, c.CommonWeight)
	}

	rows, err := q.db.QueryContext(ctx, termsSearch, query, pos, c.PosWeight, c.PopLog, c.PopWeight, c.CommonWeight, c.LenLog)
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
