// Package testdb contains test helper functions related to db
package testdb

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/lang"
	"github.com/s12chung/text2anki/pkg/util/test"
)

// MustSetupAndSeed calls Setup() and Seed(), if it fails, it exits
func MustSetupAndSeed(testName string) {
	if err := SetupTempDB(testName); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := Seed(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// SetupAndSeed calls Setup() and Seed()
func SetupAndSeed(t *testing.T, testName string) {
	SetupTempDBT(t, testName)
	SeedT(t)
}

// SetupTempDB calls db.SetDB with a temp file
func SetupTempDB(testName string) error {
	filename := test.GenerateFilename(testName, ".sqlite3")
	if err := db.SetDB(path.Join(os.TempDir(), filename)); err != nil {
		return err
	}
	return db.Create(context.Background())
}

// SetupTempDBT calls SetupTempDB and checks errors
func SetupTempDBT(t *testing.T, testName string) {
	require := require.New(t)
	err := SetupTempDB(testName)
	require.NoError(err)
}

// SeedT seeds the database with a small amount of data
func SeedT(t *testing.T) {
	require := require.New(t)
	require.NoError(Seed())
}

// Seed seeds the database with a small amount of data
func Seed() error {
	return SeedModels()
}

var modelDatas = []generateModelsCodeData{
	{Name: "Term", CreateCode: "queries.TermCreate(context.Background(), term.CreateParams())"},
	{Name: "SourceSerialized", CreateCode: "queries.SourceSerializedCreate(context.Background(), sourceSerialized.TokenizedTexts)"},
}

type generateModelsCodeData struct {
	Name       string
	CreateCode string
}

//go:embed generate_models.go.tmpl
var generateModelsCodeTemplate string

// GenerateModelsCode generates code for the testdb models
func GenerateModelsCode() ([]byte, error) {
	temp, err := template.New("top").Funcs(template.FuncMap{
		"pluralize": pluralize,
		"lower":     lower,
	}).Parse(generateModelsCodeTemplate)
	if err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}
	if err = temp.Execute(&buffer, modelDatas); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "s") {
		return s
	}
	return s + "s"
}

func lower(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := strings.ToLower(string(s[0]))
	return firstChar + s[1:]
}

// SearchTerm is a search term used for tests
const SearchTerm = "마음"

// SearchPOS is the search POS for tests
const SearchPOS = lang.PartOfSpeechVerb

// SearchConfig is the config used for test searching (so it stays constant)
var SearchConfig = db.TermsSearchConfig{
	PosWeight:    10,
	PopLog:       20,
	PopWeight:    40,
	CommonWeight: 40,
	LenLog:       2,
}
