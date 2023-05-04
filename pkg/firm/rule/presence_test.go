package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type presenceStruct struct {
	Integer int
}

func TestPresenceFieldValidator_Valid(t *testing.T) {
	fieldValidator := PresenceFieldValidator{}
	tcs := []struct {
		name     string
		value    reflect.Value
		expected error
	}{
		{name: "int", value: reflect.ValueOf(10), expected: nil},
		{name: "empty_int", value: reflect.ValueOf(0), expected: errPresenceFieldValidator},
		{name: "struct", value: reflect.ValueOf(presenceStruct{Integer: 1}), expected: nil},
		{name: "empty_struct", value: reflect.ValueOf(presenceStruct{}), expected: errPresenceFieldValidator},
		{name: "func", value: reflect.ValueOf(func() {}), expected: nil},
		{name: "pointer", value: reflect.ValueOf(&presenceStruct{}), expected: nil},
		{name: "nil", value: reflect.ValueOf(nil), expected: errPresenceFieldValidator},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, fieldValidator.Valid(tc.value))
		})
	}
}
