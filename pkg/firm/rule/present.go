package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Present checks if data is non-Zero and non-Empty
type Present struct{}

// ValidateValue returns true if the data is valid (assumes ValidateType is called)
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

// ValidateType checks whether the type is valid for the Rule -- allow all types
func (p Present) ValidateType(_ reflect.Type) *firm.RuleTypeError { return nil }

var errorMapPresent = firm.ErrorMap{"Present": firm.TemplateError{Template: "is not present"}}
