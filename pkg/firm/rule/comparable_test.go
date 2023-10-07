package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

func TestEqual_ValidateValue(t *testing.T) {
	require.Equal(t, "Equal: value is not equal to 99", Equal[int]{To: 99}.ErrorMap().Error())

	tcs := []struct {
		name string
		to   int

		data     int
		hasError bool
	}{
		{name: "equal", to: 9, data: 9},
		{name: "below", to: 9, data: 1, hasError: true},
		{name: "above", to: 9, data: 100, hasError: true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testComparableRule_ValidateValue[int](t, Equal[int]{To: tc.to}, tc.hasError, tc.data)
		})
	}
}

func TestEqual_ValidateType(t *testing.T) { testComparableRule_ValidateType[Equal[int]](t) }

func TestLess_ValidateValue(t *testing.T) {
	require.Equal(t, "Less: value is not less than 99", Less[int]{To: 99}.ErrorMap().Error())
	require.Equal(t, "LessOrEqual: value is not less than or equal to 99", Less[int]{OrEqual: true, To: 99}.ErrorMap().Error())

	tcs := []struct {
		name    string
		orEqual bool
		to      int

		data     int
		hasError bool
	}{
		{name: "less", to: 9, data: 1},
		{name: "equal", to: 9, data: 9, hasError: true},
		{name: "above", to: 9, data: 100, hasError: true},
		{name: "or_equal_equal", orEqual: true, to: 9, data: 9},
		{name: "or_equal_below", orEqual: true, to: 9, data: -1},
		{name: "or_equal_above", orEqual: true, to: 9, data: 100, hasError: true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testComparableRule_ValidateValue[int](t, Less[int]{OrEqual: tc.orEqual, To: tc.to}, tc.hasError, tc.data)
		})
	}
}

func TestLess_ValidateType(t *testing.T) { testComparableRule_ValidateType[Less[int]](t) }

func TestGreater_ValidateValue(t *testing.T) {
	require.Equal(t, "Greater: value is not greater than 99", Greater[int]{To: 99}.ErrorMap().Error())
	require.Equal(t, "GreaterOrEqual: value is not greater than or equal to 99", Greater[int]{OrEqual: true, To: 99}.ErrorMap().Error())

	tcs := []struct {
		name    string
		orEqual bool
		to      int

		data     int
		hasError bool
	}{
		{name: "greater", to: 1, data: 9},
		{name: "equal", to: 9, data: 9, hasError: true},
		{name: "below", to: 100, data: 9, hasError: true},
		{name: "or_equal_equal", orEqual: true, to: 9, data: 9},
		{name: "or_equal_above", orEqual: true, to: 9, data: 100},
		{name: "or_equal_below", orEqual: true, to: 9, data: -1, hasError: true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testComparableRule_ValidateValue[int](t, Greater[int]{OrEqual: tc.orEqual, To: tc.to}, tc.hasError, tc.data)
		})
	}
}

func TestGreater_ValidateType(t *testing.T) { testComparableRule_ValidateType[Greater[int]](t) }

//nolint:revive,stylecheck // for tests
func testComparableRule_ValidateValue[T comparable](t *testing.T, rule comparableRule[T], hasError bool, data T) {
	require := require.New(t)
	var expected firm.ErrorMap
	if hasError {
		expected = rule.ErrorMap()
	}
	require.Equal(expected, rule.ValidateValue(reflect.ValueOf(data)))
}

//nolint:revive,stylecheck // for tests
func testComparableRule_ValidateType[T comparableRule[int]](t *testing.T) {
	i := 0
	badCondition := "is not a int"

	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "matching type", data: 0},
		{name: "matching type pointer", data: &i, badCondition: badCondition},
		{name: "other type", data: "", badCondition: badCondition},
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
			var rule T
			require.Equal(err, rule.ValidateType(typ))
		})
	}
}
