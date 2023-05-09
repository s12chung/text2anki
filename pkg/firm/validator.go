package firm

import (
	"reflect"
	"strconv"
)

// NewStructValidator returns a new StructValidator
func NewStructValidator(ruleMap RuleMap) StructValidator {
	rm := map[string]*[]Rule{}
	for k, v := range ruleMap {
		rules := v
		rm[k] = &rules
	}
	return StructValidator{RuleMap: rm}
}

// StructValidator validates structs
type StructValidator struct {
	RuleMap map[string]*[]Rule
}

// Validate validates the data
func (s StructValidator) Validate(data any) Result {
	return validate(data, s.ValidateMerge)
}

// ValidateValue validates the data value
func (s StructValidator) ValidateValue(value reflect.Value) ErrorMap {
	return validateValue(value, s.ValidateMerge)
}

const structValidatorErrorKey = "StructValidator"

func structValidatorError(value reflect.Value) *TemplatedError {
	return &TemplatedError{
		TemplateFields: map[string]string{"Type": typeName(value)},
		Template:       "passed in data of type, {{.Type}}, is not a Struct",
	}
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (s StructValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	if value.Kind() != reflect.Struct {
		MergeErrorMap(key, ErrorMap{structValidatorErrorKey: structValidatorError(value)}, errorMap)
		return
	}

	for fieldName, rules := range s.RuleMap {
		field, _ := value.Type().FieldByName(fieldName)
		fieldKey := joinKeys(key, ErrorKey(field.Name))
		validateMerge(value.FieldByName(fieldName), fieldKey, errorMap, *rules)
	}
}

// NewSliceValidator returns a new SliceValidator
func NewSliceValidator(elementRules ...Rule) SliceValidator {
	return SliceValidator{
		ElementRules: elementRules,
	}
}

// SliceValidator validates slices and arrys
type SliceValidator struct {
	ElementRules []Rule
}

// Validate validates the data
func (s SliceValidator) Validate(data any) Result {
	return validate(data, s.ValidateMerge)
}

// ValidateValue validates the data value
func (s SliceValidator) ValidateValue(value reflect.Value) ErrorMap {
	return validateValue(value, s.ValidateMerge)
}

const sliceValidatorErrorKey = "SliceValidator"

func sliceValidatorError(value reflect.Value) *TemplatedError {
	return &TemplatedError{
		TemplateFields: map[string]string{"Type": typeName(value)},
		Template:       "passed in data of type, {{.Type}}, is not a Slice or Array",
	}
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (s SliceValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	if !(value.Kind() == reflect.Slice || value.Kind() == reflect.Array) {
		MergeErrorMap(key, ErrorMap{sliceValidatorErrorKey: sliceValidatorError(value)}, errorMap)
		return
	}

	for i := 0; i < value.Len(); i++ {
		indexKey := joinKeys(key, ErrorKey("["+strconv.Itoa(i)+"]"))
		validateMerge(value.Index(i), indexKey, errorMap, s.ElementRules)
	}
}

// NewValueValidator returns a ValueValidator
func NewValueValidator(rules ...Rule) ValueValidator {
	return ValueValidator{
		Rules: rules,
	}
}

// ValueValidator validates a simple value
type ValueValidator struct {
	Rules []Rule
}

// Validate validates the data
func (v ValueValidator) Validate(data any) Result {
	return validate(data, v.ValidateMerge)
}

// ValidateValue validates the data value
func (v ValueValidator) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	v.ValidateMerge(value, "", errorMap)
	return errorMap
}

// ValidateMerge validates the data value, also doing a merge with the errorMap
func (v ValueValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	validateMerge(value, key, errorMap, v.Rules)
}

type validateMergeF func(value reflect.Value, key ErrorKey, errorMap ErrorMap)

func validate(data any, validateMerge validateMergeF) Result {
	value := reflect.ValueOf(data)
	errorMap := ErrorMap{}
	validateMerge(value, typeNameKey(value), errorMap)
	return MapResult{errorMap: errorMap}
}

func validateValue(value reflect.Value, validateMerge validateMergeF) ErrorMap {
	value = indirect(value)
	errorMap := ErrorMap{}

	if !value.IsValid() {
		return errorMap
	}
	validateMerge(value, "", errorMap)
	return errorMap
}

func validateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap, rules []Rule) {
	for _, rule := range rules {
		MergeErrorMap(key, rule.ValidateValue(value), errorMap)
	}
}
