// Package csv contains CSV file helpers
package csv

import (
	"encoding/csv"
	"os"

	"github.com/s12chung/text2anki/pkg/util/ioutil"
)

// File stores the array into a CSV
func File(path string, rows [][]string) error {
	//nolint:gosec // generic library
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, ioutil.OwnerRWGroupR)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck,gosec // just closing file

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range rows {
		if err = writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
