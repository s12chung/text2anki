// Package main is the start point for db
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/s12chung/text2anki/db/pkg/cli/search"
	"github.com/s12chung/text2anki/db/pkg/csv"
	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

func init() {
	flag.Parse()
}

const cmdStringGenerate = "generate"
const cmdStringDiff = "diff"
const cmdStringCreate = "create"
const cmdStringSeed = "seed"
const cmdStringSchema = "schema"
const cmdStringSearch = "search"
const commands = cmdStringGenerate + "/" + cmdStringDiff + "/" + cmdStringCreate + "/" + cmdStringSeed + "/" + cmdStringSchema + "/" + cmdStringSchema
const usage = "usage: %v [" + commands + "]"

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
	case cmdStringGenerate:
		return cmdGenerate()
	case cmdStringDiff:
		return cmdDiff()
	case cmdStringCreate:
		return cmdCreate()
	case cmdStringSeed:
		return cmdSeed()
	case cmdStringSchema:
		return cmdSchema()
	case cmdStringSearch:
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

const generateFile = "pkg/db/testdb/models.go"

func cmdGenerate() error {
	code, err := testdb.GenerateModelsCode()
	if err != nil {
		return err
	}
	return os.WriteFile(generateFile, code, ioutil.OwnerRWGroupR)
}

func cmdDiff() error {
	existingCode, err := os.ReadFile(generateFile)
	if err != nil {
		return err
	}

	code, err := testdb.GenerateModelsCode()
	if err != nil {
		return err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(existingCode)),
		B:        difflib.SplitLines(string(code)),
		FromFile: "Original",
		ToFile:   "Generated",
		Context:  3,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return err
	}
	if text != "" {
		fmt.Println(text)
		return fmt.Errorf("diff exists")
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

const searchConfigPath = "tmp/search.json"

func cmdSearch() error {
	if err := setDB(); err != nil {
		return err
	}

	config, err := search.GetOrDefaultConfig(searchConfigPath)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(config, search.Config{}) {
		fmt.Println("Wrote search config to " + searchConfigPath + ", edit it and run command again")
		return nil
	}

	validation := firm.Validate(config)
	if !validation.IsValid() {
		return fmt.Errorf("config is missing a field: %v", validation)
	}

	for _, query := range config.Queries {
		terms, err := db.Qs().TermsSearch(context.Background(), query.Str, query.POS)
		if err != nil {
			return err
		}
		rows, err := search.TermsSearchToCSVRows(terms)
		if err != nil {
			return err
		}
		rows = append(search.ConfigToCSVRows(), rows...)

		filename := "tmp/search-" + query.Str
		if query.POS != "" {
			filename += "_" + string(query.POS) + "_"
		}
		if err = csv.File(filename+".csv", rows); err != nil {
			return err
		}
	}
	return nil
}
