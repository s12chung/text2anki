package rule

import (
	"reflect"
	"strings"

	"github.com/s12chung/text2anki/pkg/firm"
)

// TrimPresent checks if data is not "" when strings.TrimSpace is applied
type TrimPresent struct{}

// ValidateValue returns true if the data is valid (assumes ValidateType is called)
func (t TrimPresent) ValidateValue(value reflect.Value) firm.ErrorMap {
	if strings.TrimSpace(value.String()) == "" {
		return errorMapTrimPresent
	}
	return nil
}

// ValidateType checks whether the type is valid for the Rule
func (t TrimPresent) ValidateType(typ reflect.Type) *firm.RuleTypeError {
	if typ.Kind() != reflect.String {
		return firm.NewRuleTypeError(typ, "is not a String")
	}
	return nil
}

var errorMapTrimPresent = firm.ErrorMap{"TrimPresent": firm.TemplateError{Template: "is just spaces or empty"}}
