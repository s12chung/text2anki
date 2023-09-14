// Package ziputil contains utilities for archive/zip
package ziputil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipDir returns the []byte of a zipped directory
func ZipDir(dirToZip string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buffer)

	if err := filepath.Walk(dirToZip, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimLeft(strings.TrimPrefix(path, dirToZip), "/")
		if header.Name == "" {
			return nil
		}
		if info.IsDir() {
			header.Name += "/"
			_, err := zipWriter.CreateHeader(header)
			return err
		}

		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(path) //nolint:gosec // purpose of function
		if err != nil {
			return err
		}
		defer file.Close() //nolint:errcheck // ok if it fails

		if _, err = io.Copy(writer, file); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	err := zipWriter.Close() // must be called before buffer.Bytes()
	return buffer.Bytes(), err
}

// CompareContents compares the contents of the zip file to the paths and handlePath
func CompareContents(zipBytes []byte, paths []string, handlePath func(path string, contents []byte)) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return err
	}
	for i, zipFile := range zipReader.File {
		if paths[i] != zipFile.Name {
			return fmt.Errorf("non matching path (%v) at, %v: %v", paths[i], i, zipFile.Name)
		}
		if zipFile.FileInfo().IsDir() {
			continue
		}

		fileReader, err := zipFile.Open()
		if err != nil {
			return err
		}
		contents, err := io.ReadAll(fileReader)
		if err != nil {
			return err
		}
		if err := fileReader.Close(); err != nil {
			return err
		}
		handlePath(zipFile.Name, contents)
	}
	return nil
}
