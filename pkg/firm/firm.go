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
	return ErrorMap{
		notFoundRuleErrorKey: notFoundRuleError(value),
	}
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
	ValidateMerge(value reflect.Value, key string, errorMap ErrorMap)
}

// RuleMap is a map of fields or keys to rules
type RuleMap map[string][]Rule

// ErrorMap is a map of TemplatedError keys to their respective TemplatedError
type ErrorMap map[string]*TemplatedError

// MergeErrorMap merges src into dest, given appending path to the src keys
func MergeErrorMap(path string, src, dest ErrorMap) {
	for k, v := range src {
		dest[joinKeys(path, k)] = v
	}
}

// TemplatedError is an error that contains a key matching a field or top level, a golang template, and template fields
type TemplatedError struct {
	Key            string
	Template       string
	TemplateFields map[string]string
}

// Error returns a string for the error
func (t *TemplatedError) Error() string {
	keyPrefix := ""
	if t.Key != "" {
		keyPrefix = "invalid " + t.Key + ": "
	}
	return keyPrefix + t.templateError()
}

func (t *TemplatedError) templateError() string {
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
	name := nilName
	if value.IsValid() {
		value = reflect.Indirect(value)
		name = value.Type().Name()
	}
	return name
}

func joinKeys(keys ...string) string {
	keysCopy := make([]string, len(keys))
	i := 0
	for _, v := range keys {
		if v == "" {
			continue
		}
		keysCopy[i] = v
		i++
	}
	keysCopy = keysCopy[:i]
	return strings.Join(keysCopy, ".")
}
