package firm

import (
	"reflect"
)

// StructValidator validates structs
type StructValidator struct {
	Registry      *Registry
	TopLevelRules []Rule
	RuleMap       RuleMap
}

// Validate validates the data
func (s StructValidator) Validate(data any) Result {
	return MapResult{errorMap: s.ValidateValue(reflect.ValueOf(data))}
}

// ValidateValue validates the data value
func (s StructValidator) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	s.ValidateMerge(value, typeName(value), errorMap)
	return errorMap
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (s StructValidator) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	value = reflect.Indirect(value)
	validateMerge(value, key, errorMap, s.TopLevelRules)

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		rules, exists := s.RuleMap[field.Name]
		if !exists {
			continue
		}
		fieldValue := value.Field(i)
		fieldKey := joinKeys(key, field.Name)

		validateMerge(fieldValue, fieldKey, errorMap, rules)

		indirectValue := reflect.Indirect(fieldValue)
		validator := s.Registry.Validator(indirectValue)
		if validator != nil {
			validator.ValidateMerge(indirectValue, fieldKey, errorMap)
		}
	}
}

// NewValueValidator returns a ValueValidator
func NewValueValidator(rules ...Rule) ValueValidator {
	return ValueValidator{
		ValueRules: rules,
	}
}

// ValueValidator validates a simpel value
type ValueValidator struct {
	ValueRules []Rule
}

// Validate validates the data
func (v ValueValidator) Validate(data any) Result {
	return MapResult{errorMap: v.ValidateValue(reflect.ValueOf(data))}
}

// ValidateValue validates the data value
func (v ValueValidator) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	v.ValidateMerge(value, typeName(value), errorMap)
	return errorMap
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (v ValueValidator) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	value = reflect.Indirect(value)
	validateMerge(value, key, errorMap, v.ValueRules)
}

func validateMerge(value reflect.Value, key string, errorMap ErrorMap, rules []Rule) {
	for _, rule := range rules {
		MergeErrorMap(key, rule.ValidateValue(value), errorMap)
	}
}
