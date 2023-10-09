package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeName(t *testing.T) {
	i := 0

	tcs := []struct {
		name     string
		data     any
		expected string
	}{
		{name: "normal", data: i, expected: "int"},
		{name: "pointer", data: &i, expected: "int"},
		{name: "slice", data: []int{}, expected: "[]int"},
		{name: "struct", data: parent{}, expected: "firm.parent"},
		{name: "nil", data: nil, expected: "nil"},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, typeName(reflect.ValueOf(tc.data)))
		})
	}
}
