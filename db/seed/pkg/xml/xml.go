// Package xml contains XML helpers for database handling
package xml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"sort"
)

// SchemaNode represents a tyoe node in the XML file in the schema
type SchemaNode struct {
	Many         bool
	Attrs        Attrs                  `json:",omitempty"`
	NodeAttrs    Attrs                  `json:",omitempty"`
	Children     map[string]*SchemaNode `json:",omitempty"`
	childrenMany map[string]uint
}

// Attrs is a map of XML attributes marshaled into an array, map values should always be true
type Attrs map[string]bool

// MarshalJSON marshalls the Attr map into an array
func (attrs *Attrs) MarshalJSON() ([]byte, error) {
	keys := make([]string, len(*attrs))
	i := 0
	for k := range *attrs {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return json.Marshal(keys)
}

// UnmarshalJSON unmarshalls the from an array to a Attr map
func (attrs *Attrs) UnmarshalJSON(data []byte) error {
	*attrs = Attrs{}
	keys := []string{}
	if err := json.Unmarshal(data, &keys); err != nil {
		return err
	}

	for _, key := range keys {
		(*attrs)[key] = true
	}
	return nil
}

// NewSchemaNode returns a new SchemaNode
func NewSchemaNode() *SchemaNode {
	return &SchemaNode{
		Attrs:        Attrs{},
		NodeAttrs:    Attrs{},
		Children:     map[string]*SchemaNode{},
		childrenMany: map[string]uint{},
	}
}

// Schema generates a shecma from the given xml bytes, the schema is merged with current
func Schema(xmlBytes []byte, current *SchemaNode) (*SchemaNode, error) {
	parents := []*SchemaNode{}
	decoder := xml.NewDecoder(bytes.NewReader(xmlBytes))
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF && token == nil {
				break
			}
			return nil, err
		}
		switch element := token.(type) {
		case xml.StartElement:
			currentName := element.Name.Local
			setChildren(current, currentName)

			parents, current = append(parents, current), current.Children[currentName]

			setAttrs(current, element)
		case xml.EndElement:
			setMany(current)
			current, parents = parents[len(parents)-1], parents[:len(parents)-1]
		}
	}
	current.childrenMany = map[string]uint{}
	return current, nil
}

func setChildren(current *SchemaNode, currentName string) {
	if _, exists := current.Children[currentName]; !exists {
		current.Children[currentName] = NewSchemaNode()
	}

	if _, exists := current.childrenMany[currentName]; !exists {
		current.childrenMany[currentName] = 0
	} else {
		current.childrenMany[currentName]++
	}
}

func setAttrs(current *SchemaNode, element xml.StartElement) {
	nodeAttr := ""
	for _, attr := range element.Attr {
		name := attr.Name.Local
		current.Attrs[name] = true
		if name == "att" {
			nodeAttr = attr.Value
		}
	}

	if nodeAttr != "" {
		current.NodeAttrs[nodeAttr] = true
	}
}

func setMany(current *SchemaNode) {
	for k, v := range current.childrenMany {
		if v > 0 {
			current.Children[k].Many = true
		}
	}
	current.childrenMany = map[string]uint{}
}
