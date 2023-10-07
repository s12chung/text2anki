package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

func TestTrimPresent_ValidateValue(t *testing.T) {
	require.Equal(t, "TrimPresent: value is just spaces or empty", errorMapTrimPresent.Error())

	tcs := []struct {
		name     string
		value    any
		expected firm.ErrorMap
	}{
		{name: "valid", value: "\t not space \n", expected: nil},
		{name: "just space", value: "\t \t\n \n", expected: errorMapTrimPresent},
		{name: "empty", value: "", expected: errorMapTrimPresent},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, TrimPresent{}.ValidateValue(reflect.ValueOf(tc.value)))
		})
	}
}

func TestTrimPresent_ValidateType(t *testing.T) {
	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "string", data: ""},
		{name: "not string", data: 1, badCondition: "is not a string"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			typ := reflect.TypeOf(tc.data)

			var err *firm.RuleTypeError
			if tc.badCondition != "" {
				err = firm.NewRuleTypeError(typ, tc.badCondition)
			}
			require.Equal(err, TrimPresent{}.ValidateType(typ))
		})
	}
}
