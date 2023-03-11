package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v srcDir\n", os.Args[0])
		os.Exit(-1)
	}

	srcDir := os.Args[1]

	if err := run(srcDir); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run(srcDir string) error {
	srcDir = filepath.Clean(srcDir)
	destDir, err := setupDestDir(srcDir)
	if err != nil {
		return err
	}
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == srcDir {
				return nil
			}
			return fs.SkipDir
		}

		outputPath := filepath.Join(destDir, filepath.Base(path))
		if _, err = os.Stat(outputPath); !os.IsNotExist(err) {
			return nil
		}

		contents, err := syncFileContents(path)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(outputPath, []byte(strings.Join(contents, "\n")), 0600)
		if err != nil {
			return err
		}
		return nil
	})
}

func setupDestDir(srcDir string) (string, error) {
	destDir := srcDir + "_syncfiltered"
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		return destDir, nil
	}
	return destDir, os.Mkdir(destDir, 0750)
}

func syncFileContents(srcFile string) ([]string, error) {
	//nolint:gosec // generic library
	file, err := os.Open(srcFile)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck,gosec // just closing file
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fileContents := []string{}
	isSourceLanugage := true
	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "" {
			fileContents = append(fileContents, text)
			continue
		}
		if isSourceLanugage {
			fileContents = append(fileContents, text)
		}
		isSourceLanugage = !isSourceLanugage
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return fileContents, err
}
