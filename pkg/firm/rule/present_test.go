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
		data     any
		expected firm.ErrorMap
	}{
		{name: "int", data: 10, expected: nil},
		{name: "int_empty", data: 0, expected: errorMapPresent},

		{name: "string", data: "abc", expected: nil},
		{name: "string_empty", data: "", expected: errorMapPresent},

		{name: "struct", data: presentStruct{Integer: 1}, expected: nil},
		{name: "struct_pointer", data: &presentStruct{Integer: 1}, expected: nil},
		{name: "struct_empty", data: presentStruct{}, expected: errorMapPresent},
		{name: "struct_empty_pointer", data: &presentStruct{}, expected: errorMapPresent},

		{name: "slice", data: []int{1, 2}, expected: nil},
		{name: "slice_pointer", data: []int{1, 2}, expected: nil},
		{name: "slice_empty", data: []int{}, expected: errorMapPresent},
		{name: "slice_empty_pointer", data: &[]int{}, expected: errorMapPresent},

		{name: "array", data: [3]int{1, 2, 3}, expected: nil},
		{name: "array_pointer", data: &[3]int{1, 2, 3}, expected: nil},
		{name: "array_empty", data: [3]int{}, expected: errorMapPresent},
		{name: "array_empty_pointer", data: &[3]int{}, expected: errorMapPresent},

		{name: "map", data: map[int]int{1: 1, 2: 2}, expected: nil},
		{name: "map_pointer", data: &map[int]int{1: 1, 2: 2}, expected: nil},
		{name: "map_empty", data: map[int]int{}, expected: errorMapPresent},
		{name: "map_empty_pointer", data: &map[int]int{}, expected: errorMapPresent},

		{name: "func", data: func() {}, expected: nil},
		{name: "channel_empty", data: make(chan int), expected: errorMapPresent},
		{name: "nil", data: nil, expected: errorMapPresent},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, Present{}.ValidateValue(reflect.ValueOf(tc.data)))
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
