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

	definition := NewDefinition(typedDefinition{}).ValidatesTopLevel(rules...)
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

	definition := NewDefinition(typedDefinition{}).Validates(ruleMap)
	require.Equal(ruleMap, definition.RuleMap())
	require.Panics(func() {
		definition.Validates(RuleMap{})
	})
	require.Panics(func() {
		NewDefinition(typedDefinition{}).Validates(RuleMap{"DoesNotExist": {}})
	})
}
