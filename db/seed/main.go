package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/s12chung/text2anki/db/seed/pkg/cmd/krdict"
)

func init() {
	flag.Parse()
}

const usage = "usage: %v [seed/schema]"

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
	case "seed":
		return cmdSeed()
	case "schema":
		return cmdSchema()
	default:
		return fmt.Errorf(usage+" -- %v not found", os.Args[0], cmd)
	}
}

func cmdSeed() error {
	return nil
}

func cmdSchema() error {
	node, err := krdict.RscSchema()
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
