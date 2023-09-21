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
		value    any
		expected firm.ErrorMap
	}{
		{name: "int", value: 10, expected: nil},
		{name: "empty_int", value: 0, expected: errorMapPresent},
		{name: "string", value: "abc", expected: nil},
		{name: "empty_string", value: "", expected: errorMapPresent},
		{name: "struct", value: presentStruct{Integer: 1}, expected: nil},
		{name: "empty_struct", value: presentStruct{}, expected: errorMapPresent},
		{name: "func", value: func() {}, expected: nil},
		{name: "empty_channel", value: make(chan int), expected: errorMapPresent},
		{name: "pointer", value: &presentStruct{}, expected: nil},
		{name: "slice", value: []int{1, 2}, expected: nil},
		{name: "slice_pointer", value: []int{1, 2}, expected: nil},
		{name: "empty_slice", value: []int{}, expected: errorMapPresent},
		{name: "array", value: [3]int{1, 2, 3}, expected: nil},
		{name: "array_pointer", value: &[3]int{1, 2, 3}, expected: nil},
		{name: "array_empty", value: [3]int{}, expected: errorMapPresent},
		{name: "map", value: map[int]int{1: 1, 2: 2}, expected: nil},
		{name: "empty_map", value: map[int]int{}, expected: errorMapPresent},
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
