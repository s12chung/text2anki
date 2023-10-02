package firm

import (
	"reflect"
	"strconv"
)

// NewStructValidator returns a new StructValidator
func NewStructValidator(typ reflect.Type, ruleMap RuleMap) StructValidator {
	rm := map[string]*[]Rule{}
	for k, v := range ruleMap {
		rules := v
		rm[k] = &rules
	}
	return StructValidator{Type: indirectType(typ), RuleMap: rm}
}

// StructValidator validates structs
type StructValidator struct {
	Type    reflect.Type
	RuleMap map[string]*[]Rule
}

// Validate validates the data
func (s StructValidator) Validate(data any) ErrorMap { return validate(data, s) }

// ValidateValue validates the data value (assumes ValidateType is called)
func (s StructValidator) ValidateValue(value reflect.Value) ErrorMap {
	return validateValue(value, s)
}

// ValidateType checks whether the type is valid for the Rule
func (s StructValidator) ValidateType(typ reflect.Type) *RuleTypeError {
	return validateType(typ, s.Type, "Struct")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes ValidateType is called)
func (s StructValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	if !value.IsValid() {
		return
	}
	for fieldName, rules := range s.RuleMap {
		field, _ := value.Type().FieldByName(fieldName)
		fieldKey := joinKeys(key, ErrorKey(field.Name))
		validateMerge(value.FieldByName(fieldName), fieldKey, errorMap, *rules)
	}
}

// NewSliceValidator returns a new SliceValidator
func NewSliceValidator(typ reflect.Type, elementRules ...Rule) SliceValidator {
	return SliceValidator{
		Type:         indirectType(typ),
		ElementRules: elementRules,
	}
}

// SliceValidator validates slices and arrys
type SliceValidator struct {
	Type         reflect.Type
	ElementRules []Rule
}

// Validate validates the data
func (s SliceValidator) Validate(data any) ErrorMap { return validate(data, s) }

// ValidateValue validates the data value (assumes ValidateType is called)
func (s SliceValidator) ValidateValue(value reflect.Value) ErrorMap {
	return validateValue(value, s)
}

// ValidateType checks whether the type is valid for the Rule
func (s SliceValidator) ValidateType(typ reflect.Type) *RuleTypeError {
	return validateType(typ, s.Type, "Slice or Array")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes ValidateType is called)
func (s SliceValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
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
func (v ValueValidator) Validate(data any) ErrorMap { return validate(data, v) }

// ValidateValue validates the data value (assumes ValidateType is called)
func (v ValueValidator) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	v.ValidateMerge(value, "", errorMap)
	return errorMap.ToNil()
}

// ValidateType checks whether the type is valid for the Rule
func (v ValueValidator) ValidateType(typ reflect.Type) *RuleTypeError {
	for _, rule := range v.Rules {
		if err := rule.ValidateType(typ); err != nil {
			return err
		}
	}
	return nil
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes ValidateType is called)
func (v ValueValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	validateMerge(value, key, errorMap, v.Rules)
}

var errInvalidValue = ErrorMap{"Validate": &TemplatedError{Template: "value is not valid"}}

func validate(data any, validator Validator) ErrorMap {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return errInvalidValue
	}
	return validateValueResult(value, validator)
}

func validateValueResult(value reflect.Value, validator Validator) ErrorMap {
	if err := validator.ValidateType(value.Type()); err != nil {
		return ErrorMap{"ValidateType": err.TemplatedError()}
	}

	errorMap := ErrorMap{}
	validator.ValidateMerge(value, typeNameKey(value), errorMap)
	return errorMap.ToNil()
}

func validateValue(value reflect.Value, validator Validator) ErrorMap {
	errorMap := ErrorMap{}
	validator.ValidateMerge(value, "", errorMap)
	return errorMap.ToNil()
}

func validateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap, rules []Rule) {
	for _, rule := range rules {
		rule.ValidateValue(value).MergeInto(key, errorMap)
	}
}

func validateType(typ, expectedType reflect.Type, kindString string) *RuleTypeError {
	iType := indirectType(typ)
	if iType != expectedType {
		return NewRuleTypeError(typ, "is not matching "+kindString+" of type "+expectedType.String())
	}
	return nil
}
