package firm

import (
	"fmt"
	"reflect"
	"strconv"
)

// MustNewStructValidator returns a new StructValidator, panics if there is an error
func MustNewStructValidator(data any, ruleMap RuleMap) StructValidator {
	return mustNewStructValidator(reflect.TypeOf(data), ruleMap)
}

func mustNewStructValidator(typ reflect.Type, ruleMap RuleMap) StructValidator {
	return mustNewValidator(func() (StructValidator, error) { return newStructValidator(typ, ruleMap) })
}

// NewStructValidator returns a new StructValidator
func NewStructValidator(data any, ruleMap RuleMap) (StructValidator, error) {
	return newStructValidator(reflect.TypeOf(data), ruleMap)
}

func newStructValidator(typ reflect.Type, ruleMap RuleMap) (StructValidator, error) {
	if typ == nil || typ.Kind() != reflect.Struct {
		return StructValidator{}, fmt.Errorf("type is not a Struct")
	}

	for fieldName, rules := range ruleMap {
		field, found := typ.FieldByName(fieldName)
		if !found {
			return StructValidator{}, fmt.Errorf("field, %v, not found in type: %v", fieldName, typ.String())
		}
		for _, rule := range rules {
			if err := rule.ValidateType(field.Type); err != nil {
				return StructValidator{}, fmt.Errorf("field, %v, in %v: %w", fieldName, typ.String(), err)
			}
		}
	}

	rm := map[string]*[]Rule{}
	for k, v := range ruleMap {
		rules := v
		rm[k] = &rules
	}
	return StructValidator{typ: indirectType(typ), ruleMap: rm}, nil
}

// StructValidator validates structs
type StructValidator struct {
	typ     reflect.Type
	ruleMap map[string]*[]Rule
}

// Type returns the Type the Validator handles
func (s StructValidator) Type() reflect.Type { return s.typ }

// Validate validates the data
func (s StructValidator) Validate(data any) ErrorMap { return validate(s, data) }

// ValidateValue validates the data value (assumes ValidateType is called)
func (s StructValidator) ValidateValue(value reflect.Value) ErrorMap { return validateValue(s, value) }

// ValidateType checks whether the type is valid for the Rule
func (s StructValidator) ValidateType(typ reflect.Type) *RuleTypeError {
	return validateType(typ, s.typ, "Struct")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes ValidateType is called)
func (s StructValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	if !value.IsValid() {
		return
	}
	for fieldName, rules := range s.ruleMap {
		field, _ := value.Type().FieldByName(fieldName)
		fieldKey := joinKeys(key, ErrorKey(field.Name))
		validateMerge(value.FieldByName(fieldName), fieldKey, errorMap, *rules)
	}
}

// MustNewSliceValidator returns a new SliceValidator, panics if there is an error
func MustNewSliceValidator(data any, elementRules ...Rule) SliceValidator {
	return mustNewSliceValidator(reflect.TypeOf(data), elementRules...)
}

func mustNewSliceValidator(typ reflect.Type, elementRules ...Rule) SliceValidator {
	return mustNewValidator(func() (SliceValidator, error) { return newSliceValidator(typ, elementRules...) })
}

// NewSliceValidator returns a new SliceValidator
func NewSliceValidator(data any, elementRules ...Rule) (SliceValidator, error) {
	return newSliceValidator(reflect.TypeOf(data), elementRules...)
}

func newSliceValidator(typ reflect.Type, elementRules ...Rule) (SliceValidator, error) {
	if typ == nil {
		return SliceValidator{}, fmt.Errorf("type, nil, is not a Slice or Array")
	}
	kind := typ.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return SliceValidator{}, fmt.Errorf("type, %v, is not a Slice or Array", typ.String())
	}

	for _, rule := range elementRules {
		if err := rule.ValidateType(typ.Elem()); err != nil {
			return SliceValidator{}, fmt.Errorf("element type: %w", err)
		}
	}
	return SliceValidator{typ: indirectType(typ), elementRules: elementRules}, nil
}

// SliceValidator validates slices and arrys
type SliceValidator struct {
	typ          reflect.Type
	elementRules []Rule
}

// Type returns the Type the Validator handles
func (s SliceValidator) Type() reflect.Type { return s.typ }

// Validate validates the data
func (s SliceValidator) Validate(data any) ErrorMap { return validate(s, data) }

// ValidateValue validates the data value (assumes ValidateType is called)
func (s SliceValidator) ValidateValue(value reflect.Value) ErrorMap { return validateValue(s, value) }

// ValidateType checks whether the type is valid for the Rule
func (s SliceValidator) ValidateType(typ reflect.Type) *RuleTypeError {
	return validateType(typ, s.typ, "Slice or Array")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes ValidateType is called)
func (s SliceValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	for i := 0; i < value.Len(); i++ {
		indexKey := joinKeys(key, ErrorKey("["+strconv.Itoa(i)+"]"))
		validateMerge(value.Index(i), indexKey, errorMap, s.elementRules)
	}
}

// MustNewValueValidator returns a ValueValidator, panics if there is an error
func MustNewValueValidator(data any, rules ...Rule) ValueValidator {
	return mustNewValueValidator(reflect.TypeOf(data), rules...)
}

func mustNewValueValidator(typ reflect.Type, rules ...Rule) ValueValidator {
	return mustNewValidator(func() (ValueValidator, error) { return newValueValidator(typ, rules...) })
}

// NewValueValidator returns a ValueValidator, if typ is nil, typ is firm.anyType{}
func NewValueValidator(data any, rules ...Rule) (ValueValidator, error) {
	return newValueValidator(reflect.TypeOf(data), rules...)
}

// newValueValidator returns a ValueValidator, if typ is nil, typ is firm.anyType{}
func newValueValidator(typ reflect.Type, rules ...Rule) (ValueValidator, error) {
	if typ == nil {
		typ = anyTyp
	}
	for _, rule := range rules {
		if err := rule.ValidateType(typ); err != nil {
			return ValueValidator{}, err
		}
	}
	return ValueValidator{typ: typ, rules: rules}, nil
}

// ValueValidator validates a simple value
type ValueValidator struct {
	typ   reflect.Type
	rules []Rule
}

// Type returns the Type the Validator handles
func (v ValueValidator) Type() reflect.Type { return v.typ }

// Validate validates the data
func (v ValueValidator) Validate(data any) ErrorMap { return validate(v, data) }

// ValidateValue validates the data value (assumes ValidateType is called)
func (v ValueValidator) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	v.ValidateMerge(value, "", errorMap)
	return errorMap.ToNil()
}

// ValidateType checks whether the type is valid for the Rule
func (v ValueValidator) ValidateType(typ reflect.Type) *RuleTypeError {
	return validateType(typ, v.typ, "")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes ValidateType is called)
func (v ValueValidator) ValidateMerge(value reflect.Value, key ErrorKey, errorMap ErrorMap) {
	value = indirect(value)
	validateMerge(value, key, errorMap, v.rules)
}

func mustNewValidator[T any](f func() (T, error)) T {
	validator, err := f()
	if err != nil {
		panic(err.Error())
	}
	return validator
}

var errInvalidValue = ErrorMap{"Validate": &TemplatedError{Template: "value is not valid"}}

func validate(validator Validator, data any) ErrorMap {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return errInvalidValue
	}
	return validateValueResult(validator, value)
}

func validateValueResult(validator Validator, value reflect.Value) ErrorMap {
	if err := validator.ValidateType(value.Type()); err != nil {
		return ErrorMap{"ValidateType": err.TemplatedError()}
	}

	errorMap := ErrorMap{}
	validator.ValidateMerge(value, TypeNameKey(value), errorMap)
	return errorMap.ToNil()
}

func validateValue(validator Validator, value reflect.Value) ErrorMap {
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
	if expectedType == anyTyp {
		return nil
	}

	iType := indirectType(typ)
	if iType != expectedType {
		if kindString != "" {
			kindString += " "
		}
		return NewRuleTypeError(typ, "is not matching "+kindString+"of type "+expectedType.String())
	}
	return nil
}
