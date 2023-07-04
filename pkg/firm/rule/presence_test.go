package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

type presenceStruct struct {
	Integer int
}

func TestPresenceRule_Valid(t *testing.T) {
	Rule := Presence{}
	tcs := []struct {
		name     string
		value    any
		expected firm.ErrorMap
	}{
		{name: "int", value: 10, expected: nil},
		{name: "empty_int", value: 0, expected: errorMapPresence},
		{name: "string", value: "abc", expected: nil},
		{name: "empty_string", value: "", expected: errorMapPresence},
		{name: "struct", value: presenceStruct{Integer: 1}, expected: nil},
		{name: "empty_struct", value: presenceStruct{}, expected: errorMapPresence},
		{name: "func", value: func() {}, expected: nil},
		{name: "empty_channel", value: make(chan int), expected: errorMapPresence},
		{name: "pointer", value: &presenceStruct{}, expected: nil},
		{name: "slice", value: []int{1, 2}, expected: nil},
		{name: "slice_pointer", value: []int{1, 2}, expected: nil},
		{name: "empty_slice", value: []int{}, expected: errorMapPresence},
		{name: "array", value: [3]int{1, 2, 3}, expected: nil},
		{name: "array_pointer", value: &[3]int{1, 2, 3}, expected: nil},
		{name: "array_empty", value: [3]int{}, expected: errorMapPresence},
		{name: "map", value: map[int]int{1: 1, 2: 2}, expected: nil},
		{name: "empty_map", value: map[int]int{}, expected: errorMapPresence},
		{name: "nil", value: nil, expected: errorMapPresence},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, Rule.ValidateValue(reflect.ValueOf(tc.value)))
		})
	}
}
