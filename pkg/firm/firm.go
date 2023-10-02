// Package firm contains functions to do validations
package firm

import (
	"reflect"
)

// RegisterType registers the TypeDefinition to the DefaultRegistry
var RegisterType = DefaultRegistry.RegisterType

// Validate validates the data with the DefaultRegistry
var Validate = DefaultRegistry.Validate

// DefaultRegistry is the registry used for global functions
var DefaultRegistry = &Registry{}

// DefaultValidator is the validator used by registries for not found types when DefaultValidator is not defined
var DefaultValidator = NewValueValidator(NotFoundRule{})

// NotFoundRule is the rule used for not found types in the DefaultValidator
type NotFoundRule struct{}

// ValidateValue validates the value (should never be called due to ValidateType)
func (n NotFoundRule) ValidateValue(value reflect.Value) ErrorMap {
	return ErrorMap{notFoundRuleErrorKey: notFoundRuleError(value)}
}
func notFoundRuleError(value reflect.Value) *TemplatedError {
	return &TemplatedError{
		TemplateFields: map[string]string{"Type": typeName(value)},
		Template:       "type, {{.Type}}, not found in Registry",
	}
}

const notFoundRuleErrorKey = "NotFound"

// ValidateType checks whether the type is valid for the Rule
func (n NotFoundRule) ValidateType(typ reflect.Type) *RuleTypeError {
	return NewRuleTypeError(typ, "is not found in Registry")
}

// NewRuleTypeError returns a new RuleTypeError
func NewRuleTypeError(typ reflect.Type, badCondition string) *RuleTypeError {
	return &RuleTypeError{Type: typ, BadCondition: badCondition}
}

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
	Validate(data any) ErrorMap
	ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap)
}
