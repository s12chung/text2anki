package typedvalidator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

type presentRule struct{}

const presentRuleKey = "presentRule"

func templateError(errorKey firm.ErrorKey) *firm.TemplateError {
	return &firm.TemplateError{TypeName: errorKey.TypeName(), ValueName: errorKey.ValueName(), Template: "test"}
}

func (p presentRule) ValidateValue(value reflect.Value) firm.ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return firm.ErrorMap{presentRuleKey: &firm.TemplateError{Template: "test"}}
	}
	return nil
}
func (p presentRule) ValidateType(_ reflect.Type) *firm.RuleTypeError { return nil }

type onlyKindRule struct{ kind reflect.Kind }

func (o onlyKindRule) ValidateValue(_ reflect.Value) firm.ErrorMap { return nil }
func (o onlyKindRule) ValidateType(typ reflect.Type) *firm.RuleTypeError {
	if typ.Kind() != o.kind {
		return firm.NewRuleTypeError(typ, "is not "+o.kind.String())
	}
	return nil
}

type testStruct struct {
	WillV string
	NoV   string
}

func TestNewStruct(t *testing.T) {
	tcs := []struct {
		name    string
		ruleMap firm.RuleMap
		err     error
	}{
		{name: "normal", ruleMap: firm.RuleMap{"WillV": []firm.Rule{presentRule{}}}},
		{name: "bad_rule", ruleMap: firm.RuleMap{"NotAField": []firm.Rule{presentRule{}}},
			err: fmt.Errorf("field, NotAField, not found in type: typedvalidator.testStruct")},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			validator, err := NewStruct(testStruct{}, tc.ruleMap)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}

			require.NoError(err)
			expected, err := firm.NewStructValidator(testStruct{}, tc.ruleMap)
			require.NoError(err)

			require.Equal(expected, validator.StructValidator)
		})
	}
}

func TestStruct_ValidateTyped(t *testing.T) {
	errorKey := firm.ErrorKey("typedvalidator.testStruct.WillV." + presentRuleKey)

	tcs := []validateTypedTC[testStruct]{
		{name: "valid", data: testStruct{WillV: "ok"}},
		{name: "invalid", data: testStruct{NoV: "not_ok"}, result: firm.ErrorMap{errorKey: templateError(errorKey)}},
	}
	testValidateTyped(t, tcs, func() (TypedValidator[testStruct], error) {
		return NewStruct(testStruct{}, firm.RuleMap{"WillV": []firm.Rule{presentRule{}}})
	})
}

func TestNewSlice(t *testing.T) {
	intKindRule := onlyKindRule{kind: reflect.Int}

	tcs := []struct {
		name         string
		elementRules []firm.Rule
		err          error
	}{
		{name: "normal", elementRules: []firm.Rule{presentRule{}}},
		{name: "bad_rule", elementRules: []firm.Rule{intKindRule},
			err: fmt.Errorf("element type: %w", intKindRule.ValidateType(reflect.TypeOf(testStruct{})))},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			validator, err := NewSlice([]testStruct{}, tc.elementRules...)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}

			require.NoError(err)
			expected, err := firm.NewSliceValidator([]testStruct{}, tc.elementRules...)
			require.NoError(err)

			require.Equal(expected, validator.SliceValidator)
		})
	}
}

func TestSlice_ValidateTyped(t *testing.T) {
	errorKey := firm.ErrorKey("[]typedvalidator.testStruct.[0]." + presentRuleKey)
	tcs := []validateTypedTC[[]testStruct]{
		{name: "valid", data: []testStruct{{WillV: "ok"}}},
		{name: "invalid", data: []testStruct{{}},
			result: firm.ErrorMap{errorKey: templateError(errorKey)}},
	}
	testValidateTyped(t, tcs, func() (TypedValidator[[]testStruct], error) {
		return NewSlice([]testStruct{}, presentRule{})
	})
}

func TestMustNewValue(t *testing.T) {
	intKindRule := onlyKindRule{kind: reflect.Int}

	tcs := []struct {
		name  string
		rules []firm.Rule
		err   error
	}{
		{name: "normal", rules: []firm.Rule{presentRule{}}},
		{name: "bad_rule", rules: []firm.Rule{intKindRule}, err: intKindRule.ValidateType(reflect.TypeOf(testStruct{}))},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			validator, err := NewValue(testStruct{}, tc.rules...)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}

			require.NoError(err)
			expected, err := firm.NewValueValidator(testStruct{}, tc.rules...)
			require.NoError(err)

			require.Equal(expected, validator.ValueValidator)
		})
	}
}

func TestValue_ValidateTyped(t *testing.T) {
	errorKey := firm.ErrorKey("typedvalidator.testStruct." + presentRuleKey)
	tcs := []validateTypedTC[testStruct]{
		{name: "valid", data: testStruct{WillV: "ok"}},
		{name: "invalid", data: testStruct{}, result: firm.ErrorMap{errorKey: templateError(errorKey)}},
	}
	testValidateTyped(t, tcs, func() (TypedValidator[testStruct], error) {
		return NewValue(testStruct{}, presentRule{})
	})
}

type validateTypedTC[T any] struct {
	name   string
	data   T
	result firm.ErrorMap
}

func testValidateTyped[T any](t *testing.T, tcs []validateTypedTC[T], newValidator func() (TypedValidator[T], error)) {
	validator, err := newValidator()
	require.NoError(t, err)

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.result, validator.ValidateTyped(tc.data))
			require.Equal(tc.result, validator.Validate(tc.data))
		})
	}
}
