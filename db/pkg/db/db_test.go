package db_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/s12chung/text2anki/db/pkg/db/testdb"
)

func TestMain(m *testing.M) {
	testdb.MustSetupAndSeed()
	if err := textTokenizer.Setup(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	code := m.Run()
	if err := textTokenizer.Cleanup(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	os.Exit(code)
}
