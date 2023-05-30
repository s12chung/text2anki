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
		{name: "struct", value: presenceStruct{Integer: 1}, expected: nil},
		{name: "empty_struct", value: presenceStruct{}, expected: errorMapPresence},
		{name: "func", value: func() {}, expected: nil},
		{name: "pointer", value: &presenceStruct{}, expected: nil},
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
