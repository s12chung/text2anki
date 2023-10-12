package attr

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// Len is a rule.Attribute that returns the reflect.Len attribute
type Len struct{}

// Name is the name of the Attribute
func (l Len) Name() string { return "Len" }

// Type is the type of the Attribute
func (l Len) Type() reflect.Type { return intType }

// Get gets the attribute value from the value
func (l Len) Get(value reflect.Value) reflect.Value { return reflect.ValueOf(value.Len()) }

// TypeCheck checks whether the type is valid for the Attribute
func (l Len) TypeCheck(typ reflect.Type) *firm.RuleTypeError {
	//nolint:exhaustive // these are the only types that return nil
	switch typ.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan, reflect.String:
		return nil
	case reflect.Ptr:
		// reflect.Len only handles pointers to Arrays
		if typ.Elem().Kind() == reflect.Array {
			return nil
		}
	}
	return firm.NewRuleTypeError(typ, "does not have a length (not a Slice, Array, Array pointer, Channel, Map or String)")
}
