package firm

import (
	"fmt"
	"reflect"
	"strconv"
)

// MustNewStruct returns a new Struct, panics if there is an error
func MustNewStruct[T any](ruleMap RuleMap) Struct[T] {
	return mustNewValidator(func() (Struct[T], error) { return NewStruct[T](ruleMap) })
}

// NewStruct returns a new Struct
func NewStruct[T any](ruleMap RuleMap) (Struct[T], error) {
	var zero T
	structAny, err := NewStructAny(reflect.TypeOf(zero), ruleMap)
	return Struct[T]{StructAny: structAny}, err
}

// NewStructAny returns a new StructAny
func NewStructAny(typ reflect.Type, ruleMap RuleMap) (StructAny, error) {
	if typ == nil || typ.Kind() != reflect.Struct {
		return StructAny{}, fmt.Errorf("type is not a Struct")
	}

	for fieldName, rules := range ruleMap {
		field, found := typ.FieldByName(fieldName)
		if !found {
			return StructAny{}, fmt.Errorf("field, %v, not found in type: %v", fieldName, typ.String())
		}
		for _, rule := range rules {
			if err := rule.TypeCheck(field.Type); err != nil {
				return StructAny{}, fmt.Errorf("field, %v, in %v: %w", fieldName, typ.String(), err)
			}
		}
	}

	rm := map[string]*[]Rule{}
	for k, v := range ruleMap {
		rules := v
		rm[k] = &rules
	}
	return StructAny{typ: indirectType(typ), ruleMap: rm}, nil
}

// Struct validates structs
type Struct[T any] struct{ StructAny }

// ValidateX is firm.Validator(), but with a typed arg, so no type checking is done on runtime
func (s Struct[T]) ValidateX(data T) ErrorMap { return validateX(s, data) }

// StructAny is a Struct without generics
type StructAny struct {
	typ     reflect.Type
	ruleMap map[string]*[]Rule
}

// Type returns the Type the Validator handles
func (s StructAny) Type() reflect.Type { return s.typ }

// Validate validates the data
func (s StructAny) Validate(data any) ErrorMap { return validate(s, data) }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (s StructAny) ValidateValue(value reflect.Value) ErrorMap { return validateValue(s, value) }

// TypeCheck checks whether the type is valid for the Rule
func (s StructAny) TypeCheck(typ reflect.Type) *RuleTypeError {
	return typeCheck(typ, s.typ, "Struct")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (s StructAny) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	value = indirect(value)
	if !value.IsValid() {
		return
	}
	for fieldName, rules := range s.ruleMap {
		field, _ := value.Type().FieldByName(fieldName)
		validateMerge(value.FieldByName(fieldName), joinKeys(key, field.Name), errorMap, *rules)
	}
}

// RuleMap returns the rules mapped to each field
func (s StructAny) RuleMap() RuleMap {
	ruleMap := RuleMap{}
	for k, v := range s.ruleMap {
		ruleMap[k] = *v
	}
	return ruleMap
}

// MustNewSlice returns a new Slice, panics if there is an error
func MustNewSlice[T any](elementRules ...Rule) Slice[T] {
	return mustNewValidator(func() (Slice[T], error) { return NewSlice[T](elementRules...) })
}

// NewSlice returns a new Slice
func NewSlice[T any](elementRules ...Rule) (Slice[T], error) {
	var zero T
	sliceAny, err := NewSliceAny(reflect.TypeOf(zero), elementRules...)
	return Slice[T]{SliceAny: sliceAny}, err
}

// NewSliceAny returns the Slice validator without generics
func NewSliceAny(typ reflect.Type, elementRules ...Rule) (SliceAny, error) {
	if typ == nil {
		return SliceAny{}, fmt.Errorf("type, nil, is not a Slice or Array")
	}
	kind := typ.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return SliceAny{}, fmt.Errorf("type, %v, is not a Slice or Array", typ.String())
	}

	for _, rule := range elementRules {
		if err := rule.TypeCheck(typ.Elem()); err != nil {
			return SliceAny{}, fmt.Errorf("element type: %w", err)
		}
	}
	return SliceAny{typ: indirectType(typ), elementRules: elementRules}, nil
}

// Slice validates slices and arrays
type Slice[T any] struct{ SliceAny }

// ValidateX is firm.Validator(), but with a typed arg, so no type checking is done on runtime
func (s Slice[T]) ValidateX(data T) ErrorMap { return validateX(s, data) }

// SliceAny is a Slice without generics
type SliceAny struct {
	typ          reflect.Type
	elementRules []Rule
}

// Type returns the Type the Validator handles
func (s SliceAny) Type() reflect.Type { return s.typ }

// Validate validates the data
func (s SliceAny) Validate(data any) ErrorMap { return validate(s, data) }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (s SliceAny) ValidateValue(value reflect.Value) ErrorMap { return validateValue(s, value) }

// TypeCheck checks whether the type is valid for the Rule
func (s SliceAny) TypeCheck(typ reflect.Type) *RuleTypeError {
	return typeCheck(typ, s.typ, "Slice or Array")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (s SliceAny) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	value = indirect(value)
	for i := 0; i < value.Len(); i++ {
		validateMerge(value.Index(i), joinKeys(key, "["+strconv.Itoa(i)+"]"), errorMap, s.elementRules)
	}
}

// ElementRules returns the rules each element in the Slice or Array
func (s SliceAny) ElementRules() []Rule { return s.elementRules }

// MustNewValue returns a Value, panics if there is an error
func MustNewValue[T any](rules ...Rule) Value[T] {
	return mustNewValidator(func() (Value[T], error) { return NewValue[T](rules...) })
}

// NewValue returns a new Value
func NewValue[T any](rules ...Rule) (Value[T], error) {
	var zero T
	valueAny, err := NewValueAny(reflect.TypeOf(zero), rules...)
	return Value[T]{ValueAny: valueAny}, err
}

// NewValueAny returns a ValueAny
func NewValueAny(typ reflect.Type, rules ...Rule) (ValueAny, error) {
	if typ == nil {
		typ = anyTyp
	}
	for _, rule := range rules {
		if err := rule.TypeCheck(typ); err != nil {
			return ValueAny{}, err
		}
	}
	return ValueAny{typ: typ, rules: rules}, nil
}

// Value validates a simple value
type Value[T any] struct{ ValueAny }

// ValidateX is firm.Validator(), but with a typed arg, so no type checking is done on runtime
func (v Value[T]) ValidateX(data T) ErrorMap { return validateX(v, data) }

// ValueAny is a Value without generics
type ValueAny struct {
	typ   reflect.Type
	rules []Rule
}

// Type returns the Type the Validator handles
func (v ValueAny) Type() reflect.Type { return v.typ }

// Validate validates the data
func (v ValueAny) Validate(data any) ErrorMap { return validate(v, data) }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (v ValueAny) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	v.ValidateMerge(value, "", errorMap)
	return errorMap.ToNil()
}

// TypeCheck checks whether the type is valid for the Rule
func (v ValueAny) TypeCheck(typ reflect.Type) *RuleTypeError {
	return typeCheck(typ, v.typ, "")
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (v ValueAny) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	value = indirect(value)
	validateMerge(value, key, errorMap, v.rules)
}

// Rules returns the rules for ValueAny
func (v ValueAny) Rules() []Rule { return v.rules }

func mustNewValidator[T any](f func() (T, error)) T {
	validator, err := f()
	if err != nil {
		panic(err.Error())
	}
	return validator
}

var errInvalidValue = ErrorMap{"Validate": TemplateError{Template: "value is not valid"}}

func validate(validator Validator, data any) ErrorMap {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return errInvalidValue
	}
	return validateValueResult(validator, value)
}

func validateValueResult(validator Validator, value reflect.Value) ErrorMap {
	if err := validator.TypeCheck(value.Type()); err != nil {
		return ErrorMap{"TypeCheck": err.TemplateError()}
	}

	errorMap := ErrorMap{}
	validator.ValidateMerge(value, TypeName(value), errorMap)
	return errorMap.Finish()
}

func validateValue(validator Validator, value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	validator.ValidateMerge(value, "", errorMap)
	return errorMap.ToNil()
}

func validateMerge(value reflect.Value, key string, errorMap ErrorMap, rules []Rule) {
	for _, rule := range rules {
		rule.ValidateValue(value).MergeInto(key, errorMap)
	}
}

func typeCheck(typ, expectedType reflect.Type, kindString string) *RuleTypeError {
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

func validateX(validator Validator, data any) ErrorMap {
	value := reflect.ValueOf(data)
	errorMap := ErrorMap{}
	validator.ValidateMerge(value, TypeName(value), errorMap)
	return errorMap.Finish()
}
