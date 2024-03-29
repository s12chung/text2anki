// Package main is the start point for db
package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/s12chung/text2anki/db/pkg/cli/search"
	"github.com/s12chung/text2anki/db/pkg/csv"
	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/db/testdb"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

func init() {
	flag.Parse()
}

var plog = logg.Default() //nolint:forbidigo // main package

const dbPath = "data.sqlite3"

const cmdStringCreate = "create"
const cmdStringSeed = "seed"
const cmdStringTestDB = "testdb"
const cmdStringSchema = "schema"
const cmdStringSearch = "search"

var commands = strings.Join([]string{
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
		fmt.Printf(usage+"\n", os.Args[0]) //nolint:forbidigo // usage
		os.Exit(-1)
	}

	cmd := args[0]

	if err := run(cmd); err != nil {
		plog.Error("db/main", logg.Err(err))
		os.Exit(-1)
	}
}

func run(cmd string) error {
	ctx := context.Background() //nolint:forbidigo // this is main

	switch cmd {
	case cmdStringCreate:
		return cmdCreate(ctx)
	case cmdStringSeed:
		return cmdSeed(ctx)
	case cmdStringTestDB:
		return cmdTestDB(ctx)
	case cmdStringSchema:
		return cmdSchema()
	case cmdStringSearch:
		return cmdSearch(ctx)
	default:
		return fmt.Errorf(usage+" -- %v not found", os.Args[0], cmd)
	}
}

func cmdCreate(ctx context.Context) error {
	txQs, err := setDB(ctx, dbPath, db.WriteOpts())
	if err != nil {
		return err
	}
	defer txQs.Rollback()                           //nolint:errcheck // rollback can fail if committed
	if err := txQs.Create(txQs.Ctx()); err != nil { //nolint:contextcheck // this is my pattern
		return err
	}
	return txQs.Commit()
}

func cmdSeed(ctx context.Context) error {
	txQs, err := setDB(ctx, dbPath, db.WriteOpts())
	if err != nil {
		return err
	}
	defer txQs.Rollback() //nolint:errcheck // rollback can fail if committed

	if err := txQs.Create(txQs.Ctx()); err != nil { //nolint:contextcheck // this is my pattern
		return err
	}
	if err := seedkrdict.Seed(txQs, seedkrdict.DefaultRscPath); err != nil { //nolint:contextcheck // this is my pattern
		return err
	}
	if err := testdb.SeedList(txQs, map[string]bool{"Terms": false, "SourceStructureds": false}); err != nil {
		return err
	}
	return txQs.Commit()
}
func cmdTestDB(ctx context.Context) error { return testdb.Create(ctx) }

func cmdSchema() error {
	node, err := seedkrdict.RscSchema(seedkrdict.DefaultRscPath)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(bytes)) //nolint:forbidigo // it's the output of the command
	return nil
}

const searchConfigPath = "tmp/search.json"

func cmdSearch(ctx context.Context) error {
	txQs, err := setDB(ctx, dbPath, nil)
	if err != nil {
		return err
	}
	defer txQs.Rollback() //nolint:errcheck // rollback can fail if committed

	config, err := search.GetOrDefaultConfig(searchConfigPath)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(config, search.Config{}) {
		fmt.Println("Wrote search config to " + searchConfigPath + ", edit it and run command again") //nolint:forbidigo // it's the output of the command
		return nil
	}

	errorMap := firm.ValidateAny(config)
	if errorMap != nil {
		return fmt.Errorf("config is missing a field: %w", errorMap)
	}

	for _, query := range config.Queries {
		terms, err := txQs.TermsSearch(txQs.Ctx(), query.Str, query.POS) //nolint:contextcheck // this is my pattern
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

func setDB(ctx context.Context, path string, opts *sql.TxOptions) (db.TxQs, error) {
	if err := db.SetDB(path); err != nil {
		return db.TxQs{}, err
	}
	return db.NewTxQs(ctx, opts)
}
