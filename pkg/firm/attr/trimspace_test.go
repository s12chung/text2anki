package attr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrimSpace_Get(t *testing.T) {
	tcs := []struct {
		name   string
		data   any
		result string
	}{
		{name: "valid", data: "\t not space \n", result: "not space"},
		{name: "just space", data: "\t \t\n \n", result: ""},
		{name: "empty", data: "", result: ""},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.result, TrimSpace{}.Get(reflect.ValueOf(tc.data)).String())
		})
	}
}

func TestTrimSpace_TypeCheck(t *testing.T) {
	badCondition := "is not a String"

	tcs := []struct {
		name         string
		data         any
		badCondition string
	}{
		{name: "string", data: "abc"},
		{name: "not_string", data: 0, badCondition: badCondition},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testTypeCheck(t, tc.data, tc.badCondition, TrimSpace{})
		})
	}
}
