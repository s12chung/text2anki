package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Not returns Not's the firm.RuleBasic's ValidateValue()
type Not struct{ Rule firm.RuleBasic }

// ValidateValue runs all the rules after .Attr is called (assumes TypeCheck is called)
func (n Not) ValidateValue(value reflect.Value) firm.ErrorMap {
	if n.Rule.ValidateValue(value).ToNil() == nil {
		return n.ErrorMap()
	}
	return nil
}

// TypeCheck checks whether the type is valid for the Attribute
func (n Not) TypeCheck(typ reflect.Type) *firm.RuleTypeError { return n.Rule.TypeCheck(typ) }

// ErrorMap returns the ErrorMap returned from ValidateValue
func (n Not) ErrorMap() firm.ErrorMap {
	original := n.Rule.ErrorMap()
	if len(original) == 0 {
		return nil
	}

	errorMap := firm.ErrorMap{}
	for k, err := range original {
		err.Template += "--Not"
		errorMap["Not"+k] = err
	}
	return errorMap
}
