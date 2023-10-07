package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

type presentStruct struct{ Integer int }

func TestPresent_ValidateValue(t *testing.T) {
	require.Equal(t, "Present: value is not present", errorMapPresent.Error())

	tcs := []struct {
		name     string
		value    any
		expected firm.ErrorMap
	}{
		{name: "int", value: 10, expected: nil},
		{name: "int_empty", value: 0, expected: errorMapPresent},

		{name: "string", value: "abc", expected: nil},
		{name: "string_empty", value: "", expected: errorMapPresent},

		{name: "struct", value: presentStruct{Integer: 1}, expected: nil},
		{name: "struct_pointer", value: &presentStruct{Integer: 1}, expected: nil},
		{name: "struct_empty", value: presentStruct{}, expected: errorMapPresent},
		{name: "struct_empty_pointer", value: &presentStruct{}, expected: errorMapPresent},

		{name: "slice", value: []int{1, 2}, expected: nil},
		{name: "slice_pointer", value: []int{1, 2}, expected: nil},
		{name: "slice_empty", value: []int{}, expected: errorMapPresent},
		{name: "slice_empty_pointer", value: &[]int{}, expected: errorMapPresent},

		{name: "array", value: [3]int{1, 2, 3}, expected: nil},
		{name: "array_pointer", value: &[3]int{1, 2, 3}, expected: nil},
		{name: "array_empty", value: [3]int{}, expected: errorMapPresent},
		{name: "array_empty_pointer", value: &[3]int{}, expected: errorMapPresent},

		{name: "map", value: map[int]int{1: 1, 2: 2}, expected: nil},
		{name: "map_pointer", value: &map[int]int{1: 1, 2: 2}, expected: nil},
		{name: "map_empty", value: map[int]int{}, expected: errorMapPresent},
		{name: "map_empty_pointer", value: &map[int]int{}, expected: errorMapPresent},

		{name: "func", value: func() {}, expected: nil},
		{name: "channel_empty", value: make(chan int), expected: errorMapPresent},
		{name: "nil", value: nil, expected: errorMapPresent},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, Present{}.ValidateValue(reflect.ValueOf(tc.value)))
		})
	}
}

func TestPresent_ValidateType(t *testing.T) {
	type anyType struct{}

	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "any", data: anyType{}},
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
			require.Equal(err, Present{}.ValidateType(typ))
		})
	}
}
