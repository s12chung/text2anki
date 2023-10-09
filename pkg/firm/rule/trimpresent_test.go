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
		data     any
		expected firm.ErrorMap
	}{
		{name: "valid", data: "\t not space \n", expected: nil},
		{name: "just space", data: "\t \t\n \n", expected: errorMapTrimPresent},
		{name: "empty", data: "", expected: errorMapTrimPresent},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, TrimPresent{}.ValidateValue(reflect.ValueOf(tc.data)))
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
			testTypeCheck(t, tc.data, tc.badCondition, TrimPresent{})
		})
	}
}
