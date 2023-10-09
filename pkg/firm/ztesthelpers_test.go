package firm

import (
	"maps"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type onlyKindRule struct{ kind reflect.Kind }

func (o onlyKindRule) ValidateValue(_ reflect.Value) ErrorMap { return nil }
func (o onlyKindRule) TypeCheck(typ reflect.Type) *RuleTypeError {
	if typ.Kind() != o.kind {
		return NewRuleTypeError(typ, "is not "+o.kind.String())
	}
	return nil
}

type presentRule struct{}

const presentRuleKey = "presentRule"

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

func testValidates(t *testing.T, validator Validator, data any, err *TemplateError, keySuffixes ...string) {
	testValidatesFull(t, false, validator, data, err, keySuffixes...)
}

func testValidatesFull(t *testing.T, skipValidate bool, validator Validator, data any, err *TemplateError, keySuffixes ...string) {
	require := require.New(t)

	var validateValueExpected ErrorMap
	if err != nil && len(keySuffixes) > 0 {
		validateValueExpected = ErrorMap{}
		for _, key := range keySuffixes {
			validateValueExpected[ErrorKey(key)] = *err
		}
	}
	validateExpected := ErrorMap{}
	validateValueExpected.MergeInto(TypeName(reflect.ValueOf(data)), validateExpected)
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

type validateXTc[T any] struct {
	name   string
	data   T
	result ErrorMap
}

func testValidateX[T any](t *testing.T, tcs []validateXTc[T], newValidator func() (ValidatorX[T], error)) {
	validator, err := newValidator()
	require.NoError(t, err)

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.result, validator.ValidateX(tc.data))
			require.Equal(tc.result, validator.ValidateAny(tc.data))
		})
	}
}
