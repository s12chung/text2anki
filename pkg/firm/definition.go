package firm

import (
	"fmt"
	"reflect"
)

// NewDefinition returns a new Definition
func NewDefinition(data any) *Definition {
	value := indirect(reflect.ValueOf(data))
	validator := &Definition{
		typ:           value.Type(),
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
		panic(fmt.Sprintf("ValidatesTopLevel() called twice in type: %v", s.typ.Name()))
	}
	s.topLevelRules = rules
	return s
}

// Validates defines rules for fields
func (s *Definition) Validates(ruleMap RuleMap) *Definition {
	if len(s.ruleMap) != 0 {
		panic(fmt.Sprintf("Validates() called twice in type: %v", s.typ.Name()))
	}
	for fieldName := range ruleMap {
		_, exists := s.typ.FieldByName(fieldName)
		if !exists {
			panic(fmt.Sprintf("Validates() called with fieldName, %v, not in type: %v", fieldName, s.typ.Name()))
		}
	}
	s.ruleMap = ruleMap
	return s
}

// TopLevelRules return the rules that apply to the top level
func (s *Definition) TopLevelRules() []Rule {
	return s.topLevelRules
}

// RuleMap returns the map of rules for the structure
func (s *Definition) RuleMap() RuleMap {
	return s.ruleMap
}
