package firm

import (
	"fmt"
	"reflect"
)

// StructuredDefinition is a definition of a validation for structured data - maps and structs
type StructuredDefinition interface {
	Validator(registry *Registry) Validator
	TopLevelRules() []Rule
	RuleMap() RuleMap
}

// NewTypedDefinition returns a new TypedDefinition
func NewTypedDefinition(data any) *TypedDefinition {
	value := reflect.Indirect(reflect.ValueOf(data))
	validator := &TypedDefinition{
		typ:           value.Type(),
		topLevelRules: []Rule{},
		ruleMap:       RuleMap{},
	}
	return validator
}

// TypedDefinition is a definition of a validation for structs
type TypedDefinition struct {
	typ           reflect.Type
	topLevelRules []Rule
	ruleMap       RuleMap
}

// ValidatesTopLevel defines rules at top level object
func (t *TypedDefinition) ValidatesTopLevel(rules ...Rule) *TypedDefinition {
	if len(t.topLevelRules) != 0 {
		panic(fmt.Sprintf("ValidatesTopLevel() called twice in type: %v", t.typ.Name()))
	}
	t.topLevelRules = rules
	return t
}

// Validates defines rules for fields
func (t *TypedDefinition) Validates(ruleMap RuleMap) *TypedDefinition {
	if len(t.ruleMap) != 0 {
		panic(fmt.Sprintf("Validates() called twice in type: %v", t.typ.Name()))
	}
	for fieldName := range ruleMap {
		_, exists := t.typ.FieldByName(fieldName)
		if !exists {
			panic(fmt.Sprintf("Validates() called with fieldName, %v, not in type: %v", fieldName, t.typ.Name()))
		}
	}
	t.ruleMap = ruleMap
	return t
}

// Validator returns the Validator for the Definition
func (t *TypedDefinition) Validator(registry *Registry) Validator {
	return StructValidator{
		Registry:      registry,
		TopLevelRules: t.topLevelRules,
		RuleMap:       t.ruleMap,
	}
}

// TopLevelRules return the rules that apply to the top level
func (t *TypedDefinition) TopLevelRules() []Rule {
	return t.topLevelRules
}

// RuleMap returns the map of rules for the structure
func (t *TypedDefinition) RuleMap() RuleMap {
	return t.ruleMap
}
