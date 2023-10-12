package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

func TestTrimPresent_ValidateAll(t *testing.T) {
	tcs := []struct {
		name     string
		data     string
		hasError bool
	}{
		{name: "valid", data: "\t not space \n"},
		{name: "just space", data: "\t \t\n \n", hasError: true},
		{name: "empty", data: "", hasError: true},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var expected firm.ErrorMap
			if tc.hasError {
				expected = TrimPresent{}.ErrorMap()
			}
			require.Equal(t, expected, TrimPresent{}.Validate(tc.data))
			require.Equal(t, expected, TrimPresent{}.ValidateValue(reflect.ValueOf(tc.data)))
		})
	}
}

func TestTrimPresent_TypeCheck(t *testing.T) {
	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "string", data: ""},
		{name: "not string", data: 1, badCondition: "is not a String"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testTypeCheck(t, tc.data, "TrimPresent", tc.badCondition, TrimPresent{})
		})
	}
}

func TestTrimPresent_ErrorMap(t *testing.T) {
	testErrorMap(t, TrimPresent{}, "TrimPresent: value is just spaces or empty")
}
