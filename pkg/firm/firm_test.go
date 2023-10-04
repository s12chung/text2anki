package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type presentRule struct{}

const presentRuleKey = "presentRule"

var errTest = &TemplatedError{Template: "test"}

func (p presentRule) ValidateValue(value reflect.Value) ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return ErrorMap{presentRuleKey: errTest}
	}
	return nil
}
func (p presentRule) ValidateType(_ reflect.Type) *RuleTypeError { return nil }

func validateTypeErrorResult(rule Rule, data any) ErrorMap {
	return ErrorMap{"ValidateType": rule.ValidateType(reflect.TypeOf(data)).TemplatedError()}
}

func TestNotFoundRule_ValidateValue(t *testing.T) {
	require := require.New(t)
	errorMap := NotFoundRule{}.ValidateValue(reflect.ValueOf(1))
	require.NotNil(errorMap)
	require.NotEmpty(errorMap)
}
