// Package firm contains functions to do validations
package firm

import (
	"reflect"
)

type anyType struct{}

var anyTyp = AnyType()

// AnyType is a helper that returns the type used to fill in `nil` types
func AnyType() reflect.Type { return reflect.TypeOf(anyType{}) }

// MustRegisterType registers the TypeDefinition to the DefaultRegistry, panics if there is an error
var MustRegisterType = DefaultRegistry.MustRegisterType

// RegisterType registers the TypeDefinition to the DefaultRegistry
var RegisterType = DefaultRegistry.RegisterType

// Validate validates the data with the DefaultRegistry
var Validate = DefaultRegistry.Validate

// DefaultRegistry is the registry used for global functions
var DefaultRegistry = &Registry{}

// DefaultValidator is the validator used by registries for not found types when DefaultValidator is not defined
var DefaultValidator = MustNewValueValidator(nil, NotFoundRule{})

// NotFoundRule is the rule used for not found types in the DefaultValidator
type NotFoundRule struct{}

// ValidateValue validates the value
func (n NotFoundRule) ValidateValue(_ reflect.Value) ErrorMap {
	return ErrorMap{"NotFound": TemplateError{Template: "type, {{.TypeName}}, not found in Registry"}}
}

// ValidateType checks whether the type is valid for the Rule
func (n NotFoundRule) ValidateType(_ reflect.Type) *RuleTypeError { return nil }

// RuleMap is a map of fields or keys to rules
type RuleMap map[string][]Rule

// Rule defines a rule for validation definitions and validators
type Rule interface {
	ValidateValue(value reflect.Value) ErrorMap
	ValidateType(typ reflect.Type) *RuleTypeError
}

// Validator validates the data
type Validator interface {
	Rule
	Type() reflect.Type
	Validate(data any) ErrorMap
	ValidateMerge(value reflect.Value, key string, errorMap ErrorMap)
}
