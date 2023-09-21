package rule

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/firm"
)

func TestTrimPresent_ValidateValue(t *testing.T) {
	tcs := []struct {
		name     string
		value    any
		expected firm.ErrorMap
	}{
		{name: "valid", value: "\t not space \n", expected: nil},
		{name: "just space", value: "\t \t\n \n", expected: errorMapNotSpace("value is just spaces or empty")},
		{name: "empty", value: "", expected: errorMapNotSpace("value is just spaces or empty")},
		{name: "not string", value: 10, expected: errorMapNotSpace("value is not a string")},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, TrimPresent{}.ValidateValue(reflect.ValueOf(tc.value)))
		})
	}
}
