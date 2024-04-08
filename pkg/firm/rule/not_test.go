package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

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
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			var expected firm.ErrorMap
			for k, err := range tc.errorMap {
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

func TestNot_TypeCheck(t *testing.T) {
	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "string", data: ""},
		{name: "not string", data: 1, badCondition: "is not a String"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testTypeCheck(t, tc.data, "TrimPresent", tc.badCondition, Not{Rule: TrimPresent{}})
		})
	}
}

func TestNot_ErrorMap(t *testing.T) {
	rule := Not{Rule: Present{}}
	testErrorMap(t, rule, "NotPresent: value is not present--Not")
	require.Equal(t, rule.ValidateValue(reflect.ValueOf(" ")), rule.ErrorMap())
}
