package firm

import (
	"fmt"
	"reflect"
)

// NewDefinition returns a new Definition
func NewDefinition[T any]() *Definition {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Pointer {
		panic(fmt.Sprintf("NewDefinition created with pointer type, dereference it: %v", typ.String()))
	}
	validator := &Definition{
		typ:           typ,
		topLevelRules: []Rule{},
		ruleMap:       RuleMap{},
	}
	return validator
}

// Definition is a definition of a validation for structs
type Definition struct {
	typ           reflect.Type
	topLevelRules []Rule
	ruleMap       RuleMap
}

// ValidatesTopLevel defines rules at top level object
func (s *Definition) ValidatesTopLevel(rules ...Rule) *Definition {
	if len(s.topLevelRules) != 0 {
		panic(fmt.Sprintf("ValidatesTopLevel() called twice in type: %v", s.typ.String()))
	}
	s.topLevelRules = rules
	return s
}

// Validates defines rules for fields
func (s *Definition) Validates(ruleMap RuleMap) *Definition {
	if len(s.ruleMap) != 0 {
		panic(fmt.Sprintf("Validates() called twice in type: %v", s.typ.String()))
	}
	for fieldName := range ruleMap {
		_, exists := s.typ.FieldByName(fieldName)
		if !exists {
			panic(fmt.Sprintf("Validates() called with fieldName, %v, not in type: %v", fieldName, s.typ.String()))
		}
	}
	s.ruleMap = ruleMap
	return s
}

// Type returns the type for the definition
func (s *Definition) Type() reflect.Type { return s.typ }

// TopLevelRules return the rules that apply to the top level
func (s *Definition) TopLevelRules() []Rule { return s.topLevelRules }

// RuleMap returns the map of rules for the structure
func (s *Definition) RuleMap() RuleMap { return s.ruleMap }
