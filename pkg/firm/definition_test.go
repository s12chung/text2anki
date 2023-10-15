package firm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefinition_ValidatesTopLevel(t *testing.T) {
	require := require.New(t)
	rules := []Rule{presentRule{}}

	definition := NewDefinition[Child]().ValidatesTopLevel(rules...)
	require.Equal(rules, definition.TopLevelRules())
	require.Panics(func() {
		definition.ValidatesTopLevel()
	})
}

func TestDefinition_Validates(t *testing.T) {
	require := require.New(t)

	ruleMap := RuleMap{"Validates": {presentRule{}}}
	definition := NewDefinition[Child]().Validates(ruleMap)
	require.Equal(ruleMap, definition.RuleMap())

	require.Panics(func() { definition.Validates(RuleMap{}) })
	require.Panics(func() { NewDefinition[Child]().Validates(RuleMap{"DoesNotExist": {}}) })
}
