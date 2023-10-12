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
	if typ == nil {
		return StructAny{}, fmt.Errorf("type, nil, is not a Struct")
	}
	if typ.Kind() != reflect.Struct {
		return StructAny{}, fmt.Errorf("type, %v, is not a Struct", typ.String())
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
	return StructAny{typ: typ, ruleMap: rm}, nil
}

// Struct validates structs
type Struct[T any] struct{ StructAny }

// Validate is firm.Validator(), but with a typed arg, so no type checking is done on runtime
func (s Struct[T]) Validate(data T) ErrorMap { return validate(s, data) }

// StructAny is a Struct without generics
type StructAny struct {
	typ     reflect.Type
	ruleMap map[string]*[]Rule
}

// Type returns the Type the Validator handles
func (s StructAny) Type() reflect.Type { return s.typ }

// ValidateAny validates the data
func (s StructAny) ValidateAny(data any) ErrorMap { return validateAny(s, data) }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (s StructAny) ValidateValue(value reflect.Value) ErrorMap { return validateValue(s, value) }

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (s StructAny) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	if !value.IsValid() {
		return
	}
	for fieldName, rules := range s.ruleMap {
		field, _ := value.Type().FieldByName(fieldName)
		// no control over types, so indirect
		validateMerge(indirect(value.FieldByName(fieldName)), joinKeys(key, field.Name), errorMap, *rules)
	}
}

// TypeCheck checks whether the type is valid for the Rule
func (s StructAny) TypeCheck(typ reflect.Type) *RuleTypeError { return TypeCheck(typ, s.typ, "Struct") }

// RuleMap returns the rules mapped to each field
func (s StructAny) RuleMap() RuleMap {
	ruleMap := RuleMap{}
	for k, v := range s.ruleMap {
		ruleMap[k] = *v
	}
	return ruleMap
}

// MustNewSlice returns a new Slice, panics if there is an error
func MustNewSlice[T []U, U any](elementRules ...Rule) Slice[T, U] {
	return mustNewValidator(func() (Slice[T, U], error) { return NewSlice[T, U](elementRules...) })
}

// NewSlice returns a new Slice
func NewSlice[T []U, U any](elementRules ...Rule) (Slice[T, U], error) {
	var zero T
	sliceAny, err := NewSliceAny(reflect.TypeOf(zero), elementRules...)
	return Slice[T, U]{SliceAny: sliceAny}, err
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
	return SliceAny{typ: typ, elementRules: elementRules}, nil
}

// Slice validates slices and arrays
type Slice[T []U, U any] struct{ SliceAny }

// Validate is firm.Validator(), but with a typed arg, so no type checking is done on runtime
func (s Slice[T, U]) Validate(data T) ErrorMap { return validate(s, data) }

// SliceAny is a Slice without generics
type SliceAny struct {
	typ          reflect.Type
	elementRules []Rule
}

// Type returns the Type the Validator handles
func (s SliceAny) Type() reflect.Type { return s.typ }

// ValidateAny validates the data
func (s SliceAny) ValidateAny(data any) ErrorMap { return validateAny(s, data) }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (s SliceAny) ValidateValue(value reflect.Value) ErrorMap { return validateValue(s, value) }

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (s SliceAny) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	for i := 0; i < value.Len(); i++ {
		// no control over types, so indirect
		v := indirect(value.Index(i))
		validateMerge(v, joinKeys(key, "["+strconv.Itoa(i)+"]"), errorMap, s.elementRules)
	}
}

// TypeCheck checks whether the type is valid for the Rule
func (s SliceAny) TypeCheck(typ reflect.Type) *RuleTypeError {
	return TypeCheck(typ, s.typ, "Slice or Array")
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
		return ValueAny{}, fmt.Errorf("type is nil, not recommended")
	}
	if typ.Kind() == reflect.Pointer {
		return ValueAny{}, fmt.Errorf("type, %v, is a Pointer, not recommended", typ.String())
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

// Validate is firm.Validator(), but with a typed arg, so no type checking is done on runtime
func (v Value[T]) Validate(data T) ErrorMap { return validate(v, data) }

// ValueAny is a Value without generics
type ValueAny struct {
	typ   reflect.Type
	rules []Rule
}

// Type returns the Type the Validator handles
func (v ValueAny) Type() reflect.Type { return v.typ }

// ValidateAny validates the data
func (v ValueAny) ValidateAny(data any) ErrorMap { return validateAny(v, data) }

// ValidateValue validates the data value (assumes TypeCheck is called)
func (v ValueAny) ValidateValue(value reflect.Value) ErrorMap {
	errorMap := ErrorMap{}
	v.ValidateMerge(value, "", errorMap)
	return errorMap.ToNil()
}

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (v ValueAny) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	validateMerge(value, key, errorMap, v.rules)
}

// TypeCheck checks whether the type is valid for the Rule
func (v ValueAny) TypeCheck(typ reflect.Type) *RuleTypeError { return TypeCheck(typ, v.typ, "") }

// Rules returns the rules for ValueAny
func (v ValueAny) Rules() []Rule { return v.rules }

// RuleValidator is a Validator wrapper around Rule
type RuleValidator struct{ Rule }

// ValidateAny validates the data
func (r RuleValidator) ValidateAny(data any) ErrorMap { return validateAny(r, data) }

// ValidateMerge validates the data value, also doing a merge with the errorMap (assumes TypeCheck is called)
func (r RuleValidator) ValidateMerge(value reflect.Value, key string, errorMap ErrorMap) {
	validateMerge(value, key, errorMap, []Rule{r.Rule})
}

func mustNewValidator[T any](f func() (T, error)) T {
	validator, err := f()
	if err != nil {
		panic(err.Error())
	}
	return validator
}

var errInvalidValue = ErrorMap{"ValidateAny": TemplateError{Template: "value is not valid"}}

func validateAny(validator Validator, data any) ErrorMap {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return errInvalidValue
	}
	return validateValueResult(validator, value)
}

func validateValueResult(validator Validator, value reflect.Value) ErrorMap {
	// Users often don't have control over whether any is a pointer, so we're generous via indirect
	value = indirect(value)
	typ := value.Type()
	if err := validator.TypeCheck(typ); err != nil {
		return ErrorMap{"TypeCheck": err.TemplateError()}
	}

	errorMap := ErrorMap{}
	validator.ValidateMerge(value, typ.String(), errorMap)
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

func validate(validator Validator, data any) ErrorMap {
	// Users often don't have control over whether any is a pointer, so we're generous via indirect
	value := indirect(reflect.ValueOf(data))
	errorMap := ErrorMap{}
	validator.ValidateMerge(value, value.Type().String(), errorMap)
	return errorMap.Finish()
}
