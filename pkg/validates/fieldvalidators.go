package validates

import (
	"reflect"
)

// FieldValidator validates fields
type FieldValidator interface {
	Valid(value reflect.Value) error
}

var fieldValidatorMap = map[string]FieldValidator{}

// AddFieldValidator adds a FieldValidator to the fieldValidatorMap
func AddFieldValidator(name string, fieldValidator FieldValidator) {
	fieldValidatorMap[name] = fieldValidator
}
