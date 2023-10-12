package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Present checks if data is non-Zero, valid, and not of length 0
type Present struct{}

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (p Present) ValidateValue(value reflect.Value) firm.ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return errorMapPresent
	}
	//nolint:exhaustive // only checking against .Len() kinds
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map, reflect.String:
		if value.Len() == 0 {
			return errorMapPresent
		}
	case reflect.Ptr:
		return p.ValidateValue(value.Elem())
	}
	return nil
}

// TypeCheck checks whether the type is valid for the Rule -- allow all types
func (p Present) TypeCheck(_ reflect.Type) *firm.RuleTypeError { return nil }

// ErrorMap returns the ErrorMap returned from ValidateValue
func (p Present) ErrorMap() firm.ErrorMap { return errorMapPresent }

var errorMapPresent = firm.ErrorMap{"Present": firm.TemplateError{Template: "is not present"}}
