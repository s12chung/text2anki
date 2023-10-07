package rule

import (
	"fmt"
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Attribute gets the attribute value from the value
type Attribute interface {
	Name() string
	Type() reflect.Type
	Get(value reflect.Value) reflect.Value
	ValidateType(typ reflect.Type) *firm.RuleTypeError
}

// Attr is a rule that applies a firm.Rule to an Attribute value
type Attr struct {
	Of   Attribute
	Rule firm.Rule
}

// ValidateValue runs all the rules after .Attr is called (assumes ValidateType is called)
func (a Attr) ValidateValue(value reflect.Value) firm.ErrorMap {
	errorMap := firm.ErrorMap{}
	for k, v := range a.Rule.ValidateValue(a.Of.Get(value)) {
		err := v
		if err.TemplateFields == nil {
			err.TemplateFields = map[string]string{}
		}
		err.TemplateFields["AttrName"] = a.Of.Name()
		err.Template = "attribute, {{.AttrName}}, " + err.Template
		errorMap[firm.ErrorKey(a.Of.Name())+"-"+k] = err
	}
	return errorMap.ToNil()
}

// ValidateType checks whether the type is valid for the Attribute
func (a Attr) ValidateType(typ reflect.Type) *firm.RuleTypeError {
	if err := a.Of.ValidateType(typ); err != nil {
		return err
	}
	ofType := a.Of.Type()
	if err := a.Rule.ValidateType(ofType); err != nil {
		err.BadCondition = fmt.Sprintf("has Attr, %v, which ", a.Of.Name()) + err.BadCondition
		return err
	}
	return nil
}
