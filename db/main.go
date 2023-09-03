// Package main is the start point for db
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/s12chung/text2anki/db/pkg/cli/search"
	"github.com/s12chung/text2anki/db/pkg/csv"
	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/models"
	"github.com/s12chung/text2anki/db/pkg/db/testdb/testdbgen"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

func init() {
	flag.Parse()
}

const dbPath = "data.sqlite3"

const cmdStringGenerate = "generate"
const cmdStringDiff = "diff"
const cmdStringCreate = "create"
const cmdStringSeed = "seed"
const cmdStringTestDB = "testdb"
const cmdStringSchema = "schema"
const cmdStringSearch = "search"

var commands = strings.Join([]string{
	cmdStringGenerate,
	cmdStringDiff,
	cmdStringCreate,
	cmdStringSeed,
	cmdStringTestDB,
	cmdStringSchema,
	cmdStringSearch,
}, "/")
var usage = "usage: %v [" + commands + "]"

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
	case cmdStringTestDB:
		return cmdTestDB()
	case cmdStringSchema:
		return cmdSchema()
	case cmdStringSearch:
		return cmdSearch()
	default:
		return fmt.Errorf(usage+" -- %v not found", os.Args[0], cmd)
	}
}

const generateFile = "pkg/db/testdb/models/models.go"

func cmdGenerate() error {
	code, err := testdbgen.GenerateModelsCode()
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

	code, err := testdbgen.GenerateModelsCode()
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
		return fmt.Errorf("diff exists for generated result and %v", generateFile)
	}
	return nil
}

func cmdCreate() error {
	txQs, err := setDB(dbPath)
	if err != nil {
		return err
	}
	defer txQs.Rollback() //nolint:errcheck // rollback can fail if committed
	if err := txQs.Create(txQs.Ctx()); err != nil {
		return err
	}
	return txQs.Commit()
}

func cmdSeed() error {
	txQs, err := setDB(dbPath)
	if err != nil {
		return err
	}
	defer txQs.Rollback() //nolint:errcheck // rollback can fail if committed

	if err := txQs.Create(txQs.Ctx()); err != nil {
		return err
	}
	if err := seedkrdict.Seed(txQs, seedkrdict.DefaultRscPath); err != nil {
		return err
	}
	if err := models.SeedList(txQs, map[string]bool{"Terms": false}); err != nil {
		return err
	}
	return txQs.Commit()
}
func cmdTestDB() error { return testdb.Create() }

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
	txQs, err := setDB(dbPath)
	if err != nil {
		return err
	}
	defer txQs.Rollback() //nolint:errcheck // rollback can fail if committed

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
		terms, err := txQs.TermsSearch(txQs.Ctx(), query.Str, query.POS)
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
	return txQs.Commit()
}

func setDB(path string) (db.TxQs, error) {
	if err := db.SetDB(path); err != nil {
		return db.TxQs{}, err
	}
	return db.NewTxQs()
}
