package rule

import (
	"reflect"
	"strings"

	"github.com/s12chung/text2anki/pkg/firm"
)

// TrimPresent checks if data is not "" when strings.TrimSpace is applied
type TrimPresent struct{}

// ValidateValue returns true if the data is valid (assumes TypeCheck is called)
func (t TrimPresent) ValidateValue(value reflect.Value) firm.ErrorMap {
	if strings.TrimSpace(value.String()) == "" {
		return t.ErrorMap()
	}
	return nil
}

// TypeCheck checks whether the type is valid for the Rule
func (t TrimPresent) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	if typ.Kind() != reflect.String {
		return firm.NewRuleTypeError(typ, "is not a String")
	}
	return nil
}

// ErrorMap returns the ErrorMap returned from ValidateValue
func (t TrimPresent) ErrorMap() firm.ErrorMap { return errorMapTrimPresent }

var errorMapTrimPresent = firm.ErrorMap{"TrimPresent": firm.TemplateError{Template: "is just spaces or empty"}}
