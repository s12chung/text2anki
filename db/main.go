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
	"github.com/s12chung/text2anki/db/seed/pkg/cmd/krdict"
)

var seedDB = db.DefaultDB()

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
		return cmdCreateDB()
	case "seed":
		return cmdSeed()
	case "schema":
		return cmdSchema()
	default:
		return fmt.Errorf(usage+" -- %v not found", os.Args[0], cmd)
	}
}

//go:embed schema.sql
var ddl string

func cmdCreateDB() error {
	ctx := context.Background()
	if _, err := seedDB.ExecContext(ctx, ddl); err != nil {
		return err
	}
	return nil
}

func cmdSeed() error {
	if err := cmdCreateDB(); err != nil {
		return err
	}

	return krdict.Seed(seedDB, krdict.DefaultRscPath)
}

func cmdSchema() error {
	node, err := krdict.RscSchema(krdict.DefaultRscPath)
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
