package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

type presentStruct struct{ Integer int }

func TestPresent_ValidateValue(t *testing.T) {
	tcs := []struct {
		name     string
		data     any
		hasError bool
	}{
		{name: "int", data: 10},
		{name: "int_empty", data: 0, hasError: true},

		{name: "string", data: "abc"},
		{name: "string_empty", data: "", hasError: true},

		{name: "struct", data: presentStruct{Integer: 1}},
		{name: "struct_pointer", data: &presentStruct{Integer: 1}},
		{name: "struct_empty", data: presentStruct{}, hasError: true},
		{name: "struct_empty_pointer", data: &presentStruct{}, hasError: true},

		{name: "slice", data: []int{1, 2}},
		{name: "slice_pointer", data: []int{1, 2}},
		{name: "slice_empty", data: []int{}, hasError: true},
		{name: "slice_empty_pointer", data: &[]int{}, hasError: true},

		{name: "array", data: [3]int{1, 2, 3}},
		{name: "array_pointer", data: &[3]int{1, 2, 3}},
		{name: "array_empty", data: [3]int{}, hasError: true},
		{name: "array_empty_pointer", data: &[3]int{}, hasError: true},

		{name: "map", data: map[int]int{1: 1, 2: 2}},
		{name: "map_pointer", data: &map[int]int{1: 1, 2: 2}},
		{name: "map_empty", data: map[int]int{}, hasError: true},
		{name: "map_empty_pointer", data: &map[int]int{}, hasError: true},

		{name: "func", data: func() {}},
		{name: "channel_empty", data: make(chan int), hasError: true},
		{name: "nil", data: nil, hasError: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var expected firm.ErrorMap
			if tc.hasError {
				expected = Present{}.ErrorMap()
			}
			require.Equal(t, expected, Present{}.ValidateValue(reflect.ValueOf(tc.data)))
		})
	}
}

func TestPresent_TypeCheck(t *testing.T) {
	type anyType struct{}

	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "any", data: anyType{}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			testTypeCheck(t, tc.data, "Present", tc.badCondition, Present{})
		})
	}
}

func TestPresent_ErrorMap(t *testing.T) { testErrorMap(t, Present{}, "Present: value is not present") }
