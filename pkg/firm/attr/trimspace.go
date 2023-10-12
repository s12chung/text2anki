package attr

import (
	"reflect"
	"strings"

	"github.com/s12chung/text2anki/pkg/firm"
)

// TrimSpace is a rule.Attribute that calls strings.TrimSpace on the string
type TrimSpace struct{}

// Name is the name of the Attribute
func (t TrimSpace) Name() string { return "TrimSpace" }

// Type is the type of the Attribute
func (t TrimSpace) Type() reflect.Type { return stringType }

// Get gets the attribute value from the value
func (t TrimSpace) Get(value reflect.Value) reflect.Value {
	return reflect.ValueOf(strings.TrimSpace(value.String()))
}

// TypeCheck checks whether the type is valid for the Attribute
func (t TrimSpace) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	if typ.Kind() == reflect.String {
		return nil
	}
	return firm.NewRuleTypeError(t.Name(), typ, "is not a String")
}
