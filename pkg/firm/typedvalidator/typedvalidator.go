// Package typedvalidator contained typed validators with generics
package typedvalidator

import (
	"reflect"

	"github.com/s12chung/text2anki/pkg/firm"
)

// TypedValidator is a generic firm.Validator that has a typed ValidateTyped() function
type TypedValidator[T any] interface {
	firm.Validator
	ValidateTyped(data T) firm.ErrorMap
}

// MustNewStruct returns a new Struct, panics if there is an error
func MustNewStruct[T any](data T, ruleMap firm.RuleMap) Struct[T] {
	return mustNew(func() (Struct[T], error) { return NewStruct(data, ruleMap) })
}

// NewStruct returns a new Struct
func NewStruct[T any](data T, ruleMap firm.RuleMap) (Struct[T], error) {
	validator, err := firm.NewStructValidator(data, ruleMap)
	if err != nil {
		return Struct[T]{}, err
	}
	return Struct[T]{StructValidator: validator}, nil
}

// Struct is a StructValidator that's typed
type Struct[T any] struct{ firm.StructValidator }

// ValidateTyped is firm.Validator.WillV(), but with a typed arg, so no type checking is done on runtime
func (s Struct[T]) ValidateTyped(data T) firm.ErrorMap { return validateTyped(s, data) }

// MustNewSlice returns a new Slice, panics if there is an error
func MustNewSlice[T any](data T, elementRules ...firm.Rule) Slice[T] {
	return mustNew(func() (Slice[T], error) { return NewSlice[T](data, elementRules...) })
}

// NewSlice returns a new Slice
func NewSlice[T any](data T, elementRules ...firm.Rule) (Slice[T], error) {
	validator, err := firm.NewSliceValidator(data, elementRules...)
	if err != nil {
		return Slice[T]{}, err
	}
	return Slice[T]{SliceValidator: validator}, nil
}

// Slice is a SliceValidator that's typed
type Slice[T any] struct{ firm.SliceValidator }

// ValidateTyped is firm.Validator.WillV(), but with a typed arg, so no type checking is done on runtime
func (s Slice[T]) ValidateTyped(data T) firm.ErrorMap { return validateTyped(s, data) }

// MustNewValue returns a new Value, panics if there is an error
func MustNewValue[T any](data T, rules ...firm.Rule) Value[T] {
	return mustNew(func() (Value[T], error) { return NewValue[T](data, rules...) })
}

// NewValue returns a new Value
func NewValue[T any](data T, rules ...firm.Rule) (Value[T], error) {
	validator, err := firm.NewValueValidator(data, rules...)
	if err != nil {
		return Value[T]{}, err
	}
	return Value[T]{ValueValidator: validator}, nil
}

// Value is a ValueValidator that's typed
type Value[T any] struct{ firm.ValueValidator }

// ValidateTyped is firm.Validator.WillV(), but with a typed arg, so no type checking is done on runtime
func (v Value[T]) ValidateTyped(data T) firm.ErrorMap { return validateTyped(v, data) }

func mustNew[T any](f func() (T, error)) T {
	validator, err := f()
	if err != nil {
		panic(err.Error())
	}
	return validator
}

func validateTyped(validator firm.Validator, data any) firm.ErrorMap {
	value := reflect.ValueOf(data)
	errorMap := firm.ErrorMap{}
	validator.ValidateMerge(value, firm.TypeName(value), errorMap)
	return errorMap.Finish()
}
