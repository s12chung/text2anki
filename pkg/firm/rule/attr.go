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
	TypeCheck(typ reflect.Type) *firm.RuleTypeError
}

// Attr is a rule that applies a firm.Rule to an Attribute value
type Attr struct {
	Of   Attribute
	Rule firm.RuleBasic
}

// ValidateValue runs all the rules after .Attr is called (assumes TypeCheck is called)
func (a Attr) ValidateValue(value reflect.Value) firm.ErrorMap {
	return a.errorMap(a.Rule.ValidateValue(a.Of.Get(value)))
}

// TypeCheck checks whether the type is valid for the Attribute
func (a Attr) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	if err := a.Of.TypeCheck(typ); err != nil {
		return err
	}
	ofType := a.Of.Type()
	if err := a.Rule.TypeCheck(ofType); err != nil {
		err.BadCondition = fmt.Sprintf("has Attr, %v, which ", a.Of.Name()) + err.BadCondition
		return err
	}
	return nil
}

// ErrorMap returns the ErrorMap returned from ValidateValue
func (a Attr) ErrorMap() firm.ErrorMap { return a.errorMap(a.Rule.ErrorMap()) }

func (a Attr) errorMap(original firm.ErrorMap) firm.ErrorMap {
	if len(original) == 0 {
		return nil
	}

	errorMap := firm.ErrorMap{}
	for k, v := range original {
		err := v
		if err.TemplateFields == nil {
			err.TemplateFields = map[string]string{}
		}
		err.TemplateFields["AttrName"] = a.Of.Name()
		err.Template = "attribute, {{.AttrName}}, " + err.Template
		errorMap[firm.ErrorKey(a.Of.Name())+"-"+k] = err
	}
	return errorMap
}
