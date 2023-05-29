// Package testdbgen generates code for the testdb package
package testdbgen

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

var modelDatas = []generateModelsCodeData{
	{Name: "Term", CreateCode: "queries.TermCreate(context.Background(), term.CreateParams())"},
	{Name: "SourceSerialized", CreateCode: "queries.SourceCreate(context.Background(), sourceSerialized.ToSourceCreateParams())"},
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
