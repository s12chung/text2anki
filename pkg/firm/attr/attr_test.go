package attr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
	"github.com/s12chung/text2anki/pkg/firm/rule"
)

func testTypeCheck(t *testing.T, data any, ruleName, badCondition string, attr rule.Attribute) {
	require := require.New(t)

	typ := reflect.TypeOf(data)

	var ruleTypeError *firm.RuleTypeError
	if badCondition != "" {
		ruleTypeError = firm.NewRuleTypeError(ruleName, typ, badCondition)
	}
	require.Equal(ruleTypeError, attr.TypeCheck(typ))
}
