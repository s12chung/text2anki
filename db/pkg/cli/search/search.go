// Package search contains cli functions for the search cli command
package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

// TermsSearchToCSVRows gives the CSV rows for the given TermsSearchRows
func TermsSearchToCSVRows(terms []db.TermsSearchRow) ([][]string, error) {
	rows := make([][]string, len(terms)+1)
	rows[0] = []string{
		"Text", "Variants", "POS", "CommonLevel", "Explanation", "Popularity", "PosCalc", "PopCalc", "CommonCalc", "LenCalc", "RankCalc",
	}

	for i, term := range terms {
		dictTerm, err := term.Term.DictionaryTerm()
		if err != nil {
			return nil, err
		}

		rows[i+1] = []string{
			term.Text,
			term.Variants,
			term.PartOfSpeech,
			strconv.Itoa(int(term.CommonLevel)),
			dictTerm.Translations[0].Explanation,
			strconv.Itoa(int(term.Popularity)),
			fmt.Sprintf("%f", term.PosCalc.Float64),
			fmt.Sprintf("%f", term.PopCalc.Float64),
			fmt.Sprintf("%f", term.CommonCalc.Float64),
			fmt.Sprintf("%f", term.LenCalc.Float64),
			fmt.Sprintf("%f", term.PosCalc.Float64+term.PopCalc.Float64+term.CommonCalc.Float64+term.LenCalc.Float64),
		}
	}
	return rows, nil
}

// ConfigToCSVRows returns the CSV rows for the config
func ConfigToCSVRows() [][]string {
	c := db.GetTermsSearchConfig()
	return [][]string{
		{"PosWeight", "PopLog", "PopWeight", "CommonWeight", "LenLog"},
		strings.Fields(strings.Trim(fmt.Sprint([]int{c.PosWeight, c.PopLog, c.PopWeight, c.CommonWeight, c.LenLog}), "[]")),
	}
}

// Config is the Config for the search cli command
type Config struct {
	Queries []Query `json:"queries,omitempty"`
}

// Query represents a search query for the search cli command
type Query struct {
	Str string            `json:"str,omitempty"`
	POS lang.PartOfSpeech `json:"pos,omitempty"`
}

func init() {
	firm.RegisterType(firm.NewDefinition(Config{}).Validates(firm.RuleMap{
		"Queries": {},
	}))
	firm.RegisterType(firm.NewDefinition(Query{}).Validates(firm.RuleMap{
		"Str": {rule.Presence{}},
	}))
}

var defaultConfig = Config{
	Queries: []Query{
		{Str: "가", POS: "Verb"},
		{Str: "가"},
		{Str: "오"},
		{Str: "ㅂ"},
		{Str: "고 있다"},
	},
}

// GetOrDefaultConfig returns the config from given path
func GetOrDefaultConfig(p string) (Config, error) {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return Config{}, writeDefaultConfig(p)
	} else if err != nil {
		return Config{}, err
	}
	return getConfig(p)
}

func writeDefaultConfig(p string) error {
	if err := os.MkdirAll(path.Dir(p), ioutil.OwnerRWXGroupRX); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(&defaultConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, bytes, ioutil.OwnerRWGroupR)
}

func getConfig(p string) (Config, error) {
	//nolint:gosec // we just read a config file anyway
	bytes, err := os.ReadFile(p)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err = json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
