// Package firm contains functions to do validations
package firm

import (
	"reflect"
	"strings"
	"text/template"
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
type NotFoundRule struct {
}

// ValidateValue validates the value (always an error)
func (n NotFoundRule) ValidateValue(value reflect.Value) ErrorMap {
	return ErrorMap{notFoundRuleErrorKey: notFoundRuleError(value)}
}

const notFoundRuleErrorKey = "NotFound"

func notFoundRuleError(value reflect.Value) *TemplatedError {
	return &TemplatedError{
		TemplateFields: map[string]string{"Type": typeName(value)},
		Template:       "type, {{.Type}}, not found in Registry",
	}
}

// Rule defines a rule for validation definitions and validators
type Rule interface {
	ValidateValue(value reflect.Value) ErrorMap
}

// Validator validates the data
type Validator interface {
	Rule
	Validate(data any) Result
	ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap)
}

// RuleMap is a map of fields or keys to rules
type RuleMap map[string][]Rule

// ErrorMap is a map of TemplatedError keys to their respective TemplatedError
type ErrorMap map[ErrorKey]*TemplatedError

// MergeErrorMap merges src into dest, given appending path to the src keys
func MergeErrorMap(path ErrorKey, src, dest ErrorMap) {
	for k, v := range src {
		dest[joinKeys(path, k)] = v
	}
}

// TemplatedError is an error that contains a key matching a field or top level, a golang template, and template fields
type TemplatedError struct {
	Template       string
	TemplateFields map[string]string
}

// Error returns a string for the error
func (t *TemplatedError) Error() string {
	badTemplateString := t.Template + " (bad format)"
	temp, err := template.New("top").Parse(t.Template)
	if err != nil {
		return badTemplateString
	}
	var sb strings.Builder
	if err = temp.Execute(&sb, t.TemplateFields); err != nil {
		return badTemplateString
	}
	return sb.String()
}

const nilName = "nil"

func typeName(value reflect.Value) string {
	if !value.IsValid() {
		return nilName
	}
	return indirect(value).Type().Name()
}

func typeNameKey(value reflect.Value) ErrorKey {
	return ErrorKey(typeName(value))
}

func indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	return value
}

func indirectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

const keySeparator = "."

func joinKeys(keys ...ErrorKey) ErrorKey {
	var key ErrorKey
	for _, v := range keys {
		if v == "" {
			continue
		}
		if key != "" {
			key += keySeparator
		}
		key += v
	}
	return key
}

// ErrorKey is a string that has helper functions relating to error keys
type ErrorKey string

// TypeName returns the type name of the key
func (e ErrorKey) TypeName() string {
	s := string(e)
	firstIdx := strings.Index(s, keySeparator)
	if firstIdx == -1 {
		return ""
	}
	return s[:firstIdx]
}

// ErrorName returns the error name of the key
func (e ErrorKey) ErrorName() string {
	s := string(e)
	lastIdx := strings.LastIndex(s, keySeparator)
	if lastIdx == -1 {
		return ""
	}
	return s[lastIdx+len(keySeparator):]
}
