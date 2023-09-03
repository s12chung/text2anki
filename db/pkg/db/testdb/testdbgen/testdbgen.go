// Package testdbgen generates code for the testdb package
package testdbgen

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

var modelDatas = []generateModelsCodeData{
	{Name: "Term", CreateCode: "qs.TermCreate(tx.Ctx(), term.CreateParams())"},
	{Name: "SourceStructured", CreateCode: "qs.SourceCreate(tx.Ctx(), sourceStructured.CreateParams())"},
	{Name: "Note", CreateCode: "qs.NoteCreate(tx.Ctx(), note.CreateParams())"},
}

type generateModelsCodeData struct {
	Name       string
	CreateCode string
}

//go:embed generate_models.go.tmpl
var generateModelsCodeTemplate string

// GenerateModelsCode generates code for the testdb models
func GenerateModelsCode() ([]byte, error) {
	return generateModelsCodeRaw(modelDatas)
}

func generateModelsCodeRaw(modelDatas []generateModelsCodeData) ([]byte, error) {
	temp, err := template.New("top").Funcs(template.FuncMap{
		"pluralize": pluralize,
		"lower":     lower,
		"alignPad":  alignPad,
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

func alignPad(s string, pad int) string {
	return strings.Repeat(" ", pad-len(s))
}
