// Package csv contains CSV file helpers
package csv

import (
	"encoding/csv"
	"os"

	"github.com/s12chung/text2anki/pkg/util/ioutils"
)

// File stores the array into a CSV
func File(path string, rows [][]string) error {
	//nolint:gosec // generic library
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, ioutils.OwnerRWGroupR)
	if err != nil {
		return err
	}
	//nolint:errcheck,gosec // just closing file
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range rows {
		if err = writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
