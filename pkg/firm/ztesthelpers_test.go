package firm

import (
	"maps"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type presentRule struct{}

const presentRuleKey = "presentRule"

func (p presentRule) ValidateValue(value reflect.Value) ErrorMap {
	if value.IsZero() {
		return ErrorMap{presentRuleKey: *presentRuleError()}
	}
	return nil
}
func (p presentRule) ValidateType(_ reflect.Type) *RuleTypeError { return nil }

func presentRuleError() *TemplateError { return &TemplateError{Template: presentRuleKey + " template"} }

func validateTypeErrorResult(rule Rule, data any) ErrorMap {
	return ErrorMap{"ValidateType": rule.ValidateType(reflect.TypeOf(data)).TemplateError()}
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
		require.Equal(validateExpected.Finish(), validator.Validate(data))
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
