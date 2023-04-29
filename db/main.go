// Package main is the start point for db
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/db/pkg/seedkrdict"
)

func init() {
	flag.Parse()
}

const usage = "usage: %v [create/seed/schema]"

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
	default:
		return fmt.Errorf(usage+" -- %v not found", os.Args[0], cmd)
	}
}

func cmdCreate() error {
	if err := db.SetDB("data.sqlite3"); err != nil {
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
