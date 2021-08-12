package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/s12chung/text2anki/pkg/text"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v textStringFilename exportFile\n", os.Args[0])
		os.Exit(-1)
	}

	textStringFilename, exportFile := os.Args[1], os.Args[2]

	if err := run(textStringFilename, exportFile); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run(textStringFilename, exportFile string) error {
	textString, err := readTextString(textStringFilename)
	if err != nil {
		return err
	}

	parser := text.NewParser(text.Korean, text.English)
	texts, err := parser.TextsFromString(textString)
	if err != nil {
		var bytes []byte
		bytes, err = yaml.Marshal(texts)
		fmt.Println(string(bytes))
		return err
	}

	output := make([]string, len(texts))
	for i, text := range texts {
		output[i] = text.Text + "\n" + text.Translation
	}

	return ioutil.WriteFile(exportFile, []byte(strings.Join(output, "\n")), 0600)
}

func readTextString(filename string) (string, error) {
	//nolint:gosec // required for binary to work
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
