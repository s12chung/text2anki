package rule

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/attr"
)

// TrimPresent checks if data is not "" when strings.TrimSpace is applied
type TrimPresent struct{}

var trimPresentAlias = Attr{Of: attr.TrimSpace{}, Rule: Present{}}

// ValidateValue returns true if the data is valid (assumes ValidateType is called)
func (t TrimPresent) ValidateValue(value reflect.Value) firm.ErrorMap {
	if err := trimPresentAlias.ValidateValue(value); err != nil {
		return errorMapTrimPresent
	}
	return nil
}

// ValidateType checks whether the type is valid for the Rule
func (t TrimPresent) ValidateType(typ reflect.Type) *firm.RuleTypeError {
	return trimPresentAlias.ValidateType(typ)
}

var errorMapTrimPresent = firm.ErrorMap{"TrimPresent": firm.TemplateError{Template: "is just spaces or empty"}}
