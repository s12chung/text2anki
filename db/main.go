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

	"github.com/s12chung/text2anki/db/pkg/cli/search"
	"github.com/s12chung/text2anki/db/pkg/csv"
	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
	"github.com/s12chung/text2anki/pkg/validates"
)

func init() {
	flag.Parse()
}

const cmdStringCreate = "create"
const cmdStringSeed = "seed"
const cmdStringSchema = "schema"
const cmdStringSearch = "search"
const commands = cmdStringCreate + "/" + cmdStringSeed + "/" + cmdStringSchema + "/" + cmdStringSchema
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

	validator := validates.New(config)
	if !validator.IsValid() {
		return fmt.Errorf("config is missing a field at %v", validator.Key)
	}

	for _, query := range config.Queries {
		terms, err := db.Qs().TermsSearch(context.Background(), query, config.Config)
		if err != nil {
			return err
		}
		rows, err := search.TermsSearchToCSVRows(terms)
		if err != nil {
			return err
		}
		rows = append(search.ConfigToCSVRows(config), rows...)
		if err = csv.File("tmp/search-"+query+".csv", rows); err != nil {
			return err
		}
	}
	return nil
}
