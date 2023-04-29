// Package main is the start point for db
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/csv"
	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
)

func init() {
	flag.Parse()
}

const usage = "usage: %v [create/seed/schema/search]"

func main() {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf(usage+"\n", os.Args[0])
		os.Exit(-1)
	}

	cmd := args[0]

	if err := run(cmd); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run(cmd string) error {
	switch cmd {
	case "create":
		return cmdCreate()
	case "seed":
		return cmdSeed()
	case "schema":
		return cmdSchema()
	case "search":
		return cmdSearch()
	default:
		return fmt.Errorf(usage+" -- %v not found", os.Args[0], cmd)
	}
}

func setDB() error {
	if err := db.SetDB("data.sqlite3"); err != nil {
		return err
	}
	return nil
}

func cmdCreate() error {
	if err := setDB(); err != nil {
		return err
	}
	return db.Create(context.Background())
}

func cmdSeed() error {
	if err := cmdCreate(); err != nil {
		return err
	}

	return seedkrdict.Seed(context.Background(), seedkrdict.DefaultRscPath)
}

func cmdSchema() error {
	node, err := seedkrdict.RscSchema(seedkrdict.DefaultRscPath)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(bytes))
	return nil
}

var queriesString = "가,오"
var config = db.DefaultTermsSearchConfig()

func cmdSearch() error {
	if err := setDB(); err != nil {
		return err
	}

	for _, query := range strings.Split(queriesString, ",") {
		query = strings.TrimSpace(query)
		terms, err := db.Qs().TermsSearchRaw(context.Background(), query, config)
		if err != nil {
			return err
		}
		rows, err := termsSearchToCSVRows(terms)
		if err != nil {
			return err
		}
		if err = csv.File("../tmp/db-search/"+query+".csv", rows); err != nil {
			return err
		}
	}
	return nil
}

func termsSearchToCSVRows(terms []db.TermsSearchRow) ([][]string, error) {
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
