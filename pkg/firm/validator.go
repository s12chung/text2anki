package firm

import (
	"reflect"
	"strconv"
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
	s.ValidateMerge(value, typeNameKey(value), errorMap)
	return errorMap
}

const structValidatorErrorKey = "StructValidator"

func structValidatorError(value reflect.Value) *TemplatedError {
	return &TemplatedError{
		TemplateFields: map[string]string{"Type": typeName(value)},
		Template:       "passed in data of type, {{.Type}}, is not a struct",
	}
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (s StructValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	if value.Type().Kind() != reflect.Struct {
		MergeErrorMap(key, ErrorMap{
			structValidatorErrorKey: structValidatorError(value),
		}, errorMap)
		return
	}

	validateMerge(value, key, errorMap, s.TopLevelRules)

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		rules, exists := s.RuleMap[field.Name]
		if !exists {
			continue
		}
		fieldValue := value.Field(i)
		fieldKey := joinKeys(key, ErrorKey(field.Name))

		validateMerge(fieldValue, fieldKey, errorMap, rules)
		s.validateMergeRecursive(fieldValue, fieldKey, errorMap)
	}
}

func (s StructValidator) validateMergeRecursive(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	indirectValue := indirect(value)
	if indirectValue.Kind() == reflect.Array || indirectValue.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			indexKey := string(key) + "[" + strconv.Itoa(i) + "]"
			s.validateMergeRecursive(indirectValue.Index(i), ErrorKey(indexKey), errorMap)
		}
		return
	}
	validator := s.Registry.Validator(indirectValue)
	if validator != nil {
		validator.ValidateMerge(indirectValue, key, errorMap)
	}
}

// NewValueValidator returns a ValueValidator
func NewValueValidator(rules ...Rule) ValueValidator {
	return ValueValidator{
		ValueRules: rules,
	}
}

// ValueValidator validates a simple value
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
	v.ValidateMerge(value, typeNameKey(value), errorMap)
	return errorMap
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (v ValueValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	validateMerge(value, key, errorMap, v.ValueRules)
}

func validateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap, rules []Rule) {
	for _, rule := range rules {
		MergeErrorMap(key, rule.ValidateValue(value), errorMap)
	}
}
