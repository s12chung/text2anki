package firm

import (
	"maps"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const nilName = "nil"
const presentRuleKey = "presentRule"

func typeName(value reflect.Value) string {
	if !value.IsValid() {
		return nilName
	}
	return indirect(value).Type().String()
}

type onlyKindRule struct{ kind reflect.Kind }

func (o onlyKindRule) ValidateValue(_ reflect.Value) ErrorMap { return nil }
func (o onlyKindRule) TypeCheck(typ reflect.Type) *RuleTypeError {
	if typ.Kind() != o.kind {
		return NewRuleTypeError(typ, "is not "+o.kind.String())
	}
	return nil
}

type presentRule struct{}

func (p presentRule) ValidateValue(value reflect.Value) ErrorMap {
	if value.IsZero() {
		return ErrorMap{presentRuleKey: *presentRuleError("")}
	}
	return nil
}
func (p presentRule) TypeCheck(_ reflect.Type) *RuleTypeError { return nil }

//nolint:unparam // leave it for tests
func presentRuleError(errorKey ErrorKey) *TemplateError {
	return &TemplateError{ErrorKey: errorKey, Template: presentRuleKey + " template"}
}

func typeCheckErrorResult(rule Rule, data any) ErrorMap {
	return ErrorMap{"TypeCheck": rule.TypeCheck(reflect.TypeOf(data)).TemplateError()}
}

func testValidateAll(t *testing.T, validator Validator, data any, err *TemplateError, keySuffixes ...string) {
	testValidateAllFull(t, false, validator, data, err, keySuffixes...)
}

func testValidateAllFull(t *testing.T, skipValidate bool, validator Validator, data any, err *TemplateError, keySuffixes ...string) {
	require := require.New(t)

	var validateValueExpected ErrorMap
	if err != nil && len(keySuffixes) > 0 {
		validateValueExpected = ErrorMap{}
		for _, key := range keySuffixes {
			validateValueExpected[ErrorKey(key)] = *err
		}
	}
	validateExpected := ErrorMap{}
	validateValueExpected.MergeInto(typeName(reflect.ValueOf(data)), validateExpected)
	validateExpected = validateExpected.ToNil()

	if !skipValidate {
		require.Equal(validateExpected.Finish(), validator.ValidateAny(data))
	}
	require.Equal(validateValueExpected, validator.ValidateValue(reflect.ValueOf(data)))

	errorKey := "pkger.Mover.Parent"
	errorMap := ErrorMap{"Existing": TemplateError{}}
	expectedErrorMap := maps.Clone(errorMap)
	if err != nil {
		for _, keySuffix := range keySuffixes {
			expectedErrorMap[ErrorKey(joinKeys(errorKey, keySuffix))] = *err
		}
	}
	validator.ValidateMerge(reflect.ValueOf(data), errorKey, errorMap)
	require.Equal(expectedErrorMap, errorMap)
}

type validateTC[T any] struct {
	name   string
	data   T
	result ErrorMap
}

func testValidate[T any](t *testing.T, tcs []validateTC[T], newValidator func() (ValidatorTyped[T], error)) {
	validator, err := newValidator()
	require.NoError(t, err)

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.result, validator.Validate(tc.data))
			require.Equal(tc.result, validator.ValidateAny(tc.data))
		})
	}
}

func testTypeCheck(t *testing.T, data any, badCondition string, newValidator func() (Rule, error)) {
	require := require.New(t)

	validator, err := newValidator()
	require.NoError(err)

	typ := reflect.TypeOf(data)

	var ruleTypeError *RuleTypeError
	if badCondition != "" {
		ruleTypeError = NewRuleTypeError(typ, badCondition)
	}
	require.Equal(ruleTypeError, validator.TypeCheck(typ))
}
