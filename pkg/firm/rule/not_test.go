package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

func TestNot_ErrorMap(t *testing.T) {
	rule := Not{Rule: Present{}}
	testErrorMap(t, rule, "NotPresent: value is not present--Not")
	require.Equal(t, rule.ValidateValue(reflect.ValueOf(" ")), rule.ErrorMap())
}

func TestNot_ValidateValue(t *testing.T) {
	tcs := []struct {
		name string
		attr Attribute
		rule firm.RuleBasic

		data     any
		errorMap firm.ErrorMap
	}{
		{name: "normal", data: "", rule: Present{}},
		{name: "invalid", data: " ", rule: Present{}, errorMap: Present{}.ErrorMap()},
		{name: "invalid_with_empty_template_fields", data: " ", rule: Present{}, errorMap: errorMapPresent},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			var expected firm.ErrorMap
			for k, v := range tc.errorMap {
				err := v
				if expected == nil {
					expected = firm.ErrorMap{}
				}
				err.Template += "--Not"
				expected["Not"+k] = err
			}
			require.Equal(expected, Not{Rule: tc.rule}.ValidateValue(reflect.ValueOf(tc.data)))
		})
	}
}
