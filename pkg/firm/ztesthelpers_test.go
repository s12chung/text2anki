package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type presentRule struct{}

const presentRuleKey = "presentRule"

func templateError() *TemplateError { return &TemplateError{Template: "test"} }

func (p presentRule) ValidateValue(value reflect.Value) ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return ErrorMap{presentRuleKey: templateError()}
	}
	return nil
}
func (p presentRule) ValidateType(_ reflect.Type) *RuleTypeError { return nil }

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
			dupErr := *err
			validateValueExpected[ErrorKey(key)] = &dupErr
		}
	}
	validateExpected := ErrorMap{}
	validateValueExpected.MergeInto(TypeName(reflect.ValueOf(data)), validateExpected)
	validateExpected = validateExpected.ToNil()

	require.Equal(validateValueExpected, validator.ValidateValue(reflect.ValueOf(data)))

	if !skipValidate {
		// Run this last because .Finish() mutates
		validateExpected = validateExpected.Finish()
		require.Equal(validateExpected, validator.Validate(data))
	}

	// ValidateMerge will set TypeName and ValueName
	errorMap := ErrorMap{"Existing": nil}
	expectedErrorMap := ErrorMap{"Existing": nil}
	if err != nil {
		errorKey := ErrorKey("pkger.Mover.Parent") // this ErrorKey is incomplete - built via MergeInto
		for _, keySuffix := range keySuffixes {
			errorKeySuffix := ErrorKey(keySuffix)
			dupErr := *err
			expectedErrorMap[joinKeys(errorKey, errorKeySuffix)] = &dupErr
		}
	}
	validator.ValidateMerge(reflect.ValueOf(data), "pkger.Mover.Parent", errorMap)
	require.Equal(expectedErrorMap, errorMap)
}
