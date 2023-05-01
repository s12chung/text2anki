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
	"github.com/s12chung/text2anki/pkg/util/ioutils"
)

// TermsSearchToCSVRows gives the CSV rows for the given TermsSearchRows
func TermsSearchToCSVRows(terms []db.TermsSearchRow) ([][]string, error) {
	rows := make([][]string, len(terms)+1)
	rows[0] = []string{
		"Text", "Variants", "CommonLevel", "Explanation", "Popularity", "PopCalc", "CommonCalc", "LenCalc", "RankCalc",
	}

	for i, term := range terms {
		dictTerm, err := term.Term.DictionaryTerm()
		if err != nil {
			return nil, err
		}

		rows[i+1] = []string{
			term.Text,
			term.Variants,
			strconv.Itoa(int(term.CommonLevel)),
			dictTerm.Translations[0].Explanation,
			strconv.Itoa(int(term.Popularity)),
			fmt.Sprintf("%f", term.PopCalc.Float64),
			fmt.Sprintf("%f", term.CommonCalc.Float64),
			fmt.Sprintf("%f", term.LenCalc.Float64),
			fmt.Sprintf("%f", term.PopCalc.Float64+term.CommonCalc.Float64+term.LenCalc.Float64),
		}
	}
	return rows, nil
}

// ConfigToCSVRows returns the CSV rows for the config
func ConfigToCSVRows(config Config) [][]string {
	c := config.Config
	return [][]string{
		{"PopLog", "PopWeight", "CommonWeight", "LenLog"},
		strings.Fields(strings.Trim(fmt.Sprint([]int{c.PopLog, c.PopWeight, c.CommonWeight, c.LenLog}), "[]")),
	}
}

// Config is the Config for the search cli command
type Config struct {
	Queries []string             `json:"queries,omitempty" validates:"presence"`
	Config  db.TermsSearchConfig `json:"config" validates:"presence"`
}

var defaultConfig = Config{
	Queries: []string{"가", "오"},
	Config:  db.DefaultTermsSearchConfig(),
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
	if err := os.MkdirAll(path.Dir(p), ioutils.OwnerRWXGroupRX); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(&defaultConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, bytes, ioutils.OwnerRWGroupR)
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
