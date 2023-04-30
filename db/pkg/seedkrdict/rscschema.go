package seedkrdict

import (
	"os"

	"github.com/s12chung/text2anki/db/pkg/xml"
)

// RscSchema returns the krdict resource XML schema
func RscSchema(rscPath string) (*xml.SchemaNode, error) {
	xmlPaths, err := RscXMLPaths(rscPath)
	if err != nil {
		return nil, err
	}

	current := xml.NewSchemaNode()
	for _, path := range xmlPaths {
		//nolint:gosec // just parsing XML
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		current, err = xml.Schema(bytes, current)
		if err != nil {
			return nil, err
		}
	}

	return current, nil
}
