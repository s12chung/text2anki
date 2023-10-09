// Package firm contains functions to do validations
package firm

import (
	"reflect"
)

// Any is a helper that returns the type used to represent `any` in types
type Any struct{}

var anyTyp = reflect.TypeOf(Any{})

// MustRegisterType registers the TypeDefinition to the DefaultRegistry, panics if there is an error
var MustRegisterType = DefaultRegistry.MustRegisterType

// RegisterType registers the TypeDefinition to the DefaultRegistry
var RegisterType = DefaultRegistry.RegisterType

// ValidateAny validates the data with the DefaultRegistry
var ValidateAny = DefaultRegistry.ValidateAny

// DefaultRegistry is the registry used for global functions
var DefaultRegistry = &Registry{}

// DefaultValidator is the validator used by registries for not found types when DefaultValidator is not defined
var DefaultValidator = MustNewValue[Any](NotFoundRule{})

// NotFoundRule is the rule used for not found types in the DefaultValidator
type NotFoundRule struct{}

// ValidateValue validates the value
func (n NotFoundRule) ValidateValue(_ reflect.Value) ErrorMap {
	return ErrorMap{"NotFound": TemplateError{Template: "type, {{.RootTypeName}}, not found in Registry"}}
}

// TypeCheck checks whether the type is valid for the Rule
func (n NotFoundRule) TypeCheck(_ reflect.Type) *RuleTypeError { return nil }

// RuleMap is a map of fields or keys to rules
type RuleMap map[string][]Rule

// Rule defines a rule for validation definitions and validators
type Rule interface {
	ValidateValue(value reflect.Value) ErrorMap
	TypeCheck(typ reflect.Type) *RuleTypeError
}

// BasicRule is a Rule that is not composed of other rules
type BasicRule interface {
	Rule
	ErrorMap() ErrorMap
}

// Validator validates the data
type Validator interface {
	Rule
	Type() reflect.Type
	ValidateAny(data any) ErrorMap
	ValidateMerge(value reflect.Value, key string, errorMap ErrorMap)
}

// ValidatorX is a generic firm.Validator that has a typed ValidateX() function
type ValidatorX[T any] interface {
	Validator
	ValidateX(data T) ErrorMap
}
