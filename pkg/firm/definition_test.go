package firm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type typedDefinition struct {
	Primitive int
}

func TestTypedDefinition_ValidatesTopLevel(t *testing.T) {
	require := require.New(t)
	rules := []Rule{testPresence{}}

	definition := NewTypedDefinition(typedDefinition{}).ValidatesTopLevel(rules...)
	require.Equal(rules, definition.TopLevelRules())
	require.Panics(func() {
		definition.ValidatesTopLevel()
	})
}

func TestTypedDefinition_Validates(t *testing.T) {
	require := require.New(t)
	ruleMap := RuleMap{
		"Primitive": {testPresence{}},
	}

	definition := NewTypedDefinition(typedDefinition{}).Validates(ruleMap)
	require.Equal(ruleMap, definition.RuleMap())
	require.Panics(func() {
		definition.Validates(RuleMap{})
	})
	require.Panics(func() {
		NewTypedDefinition(typedDefinition{}).Validates(RuleMap{"DoesNotExist": {}})
	})
}
