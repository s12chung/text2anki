package attr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLen_Get(t *testing.T) {
	tcs := []struct {
		name   string
		data   any
		result int
	}{
		{name: "slice", data: []int{1, 2}, result: 2},
		{name: "array", data: [3]int{1, 2, 3}, result: 3},
		{name: "array_pointer", data: &[3]int{1, 2, 3}, result: 3},
		{name: "map", data: map[int]int{1: 1, 2: 2}, result: 2},
		{name: "channel", data: make(chan int), result: 0},
		{name: "string", data: "abc", result: 3},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(reflect.ValueOf(tc.result), Len{}.Get(reflect.ValueOf(tc.data)))
		})
	}
}

func TestLen_TypeCheck(t *testing.T) {
	badCondition := "does not have a length (not a Slice, Array, Array pointer, Channel, Map or String)"

	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "slice", data: []int{1, 2}},
		{name: "slice_pointer", data: &[]int{1, 2}, badCondition: badCondition},
		{name: "array", data: [3]int{1, 2, 3}},
		{name: "array_pointer", data: &[3]int{1, 2, 3}},
		{name: "map", data: map[int]int{1: 1, 2: 2}},
		{name: "map_pointer", data: &map[int]int{1: 1, 2: 2}, badCondition: badCondition},
		{name: "channel", data: make(chan int)},
		{name: "string", data: "abc"},
		{name: "int", data: 0, badCondition: badCondition},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testTypeCheck(t, tc.data, tc.badCondition, Len{})
		})
	}
}
