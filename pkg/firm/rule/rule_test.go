package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

func intEqual(i int) Equal[int] { return Equal[int]{To: i} }

func testErrorMap(t *testing.T, rule firm.RuleBasic, expected string) {
	require.Equal(t, expected, rule.ErrorMap().Error())
}

func testTypeCheck(t *testing.T, data any, badCondition string, rule firm.Rule) {
	require := require.New(t)

	typ := reflect.TypeOf(data)

	var ruleTypeError *firm.RuleTypeError
	if badCondition != "" {
		ruleTypeError = firm.NewRuleTypeError(typ, badCondition)
	}
	require.Equal(ruleTypeError, rule.TypeCheck(typ))
}
