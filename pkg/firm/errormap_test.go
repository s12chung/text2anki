package firm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorMap_Error(t *testing.T) {
	require := require.New(t)

	errorMap := ErrorMap{
		"A": &TemplatedError{Template: "field A message"},
		"B": &TemplatedError{Template: "field B message"},
	}
	require.Equal("A: field A message, B: field B message", errorMap.Error())
}

func TestErrorMap_MergeInto(t *testing.T) {
	require := require.New(t)

	src := ErrorMap{
		"A": &TemplatedError{Template: "a1"},
		"B": &TemplatedError{Template: "b1"},
	}
	dest := ErrorMap{
		"PATH.B": &TemplatedError{Template: "b2"},
		"PATH.C": &TemplatedError{Template: "c2"},
	}
	src.MergeInto("PATH", dest)
	require.Equal(ErrorMap{
		"PATH.A": &TemplatedError{Template: "a1"},
		"PATH.B": &TemplatedError{Template: "b1"},
		"PATH.C": &TemplatedError{Template: "c2"},
	}, dest)
}

func TestErrorMap_ToNil(t *testing.T) {
	tcs := []struct {
		name     string
		errorMap ErrorMap
		isNil    bool
	}{
		{name: "not_empty", errorMap: ErrorMap{"testy": nil}},
		{name: "empty", errorMap: ErrorMap{}, isNil: true},
		{name: "nil", errorMap: nil, isNil: true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			expected := tc.errorMap
			if tc.isNil {
				expected = nil
			}
			require.Equal(expected, tc.errorMap.ToNil())
		})
	}
}

func TestTemplatedError_Error(t *testing.T) {
	require := require.New(t)

	err := &TemplatedError{
		Template: "the error we go with {{ .Him }} and {{ .Her }}. yay",
		TemplateFields: map[string]string{
			"Him": "Jack",
			"Her": "Jill",
		},
	}
	require.Equal("the error we go with Jack and Jill. yay", err.Error())

	err.TemplateFields = map[string]string{
		"Her": "Jill",
	}
	require.Equal("the error we go with <no value> and Jill. yay", err.Error())

	err.Template = "{{ a }}"
	require.Equal("{{ a }} (bad format)", err.Error())
}

func TestErrorKey_TypeName(t *testing.T) {
	tcs := []struct {
		name     string
		errorKey ErrorKey
		expected string
	}{
		{name: "deep", errorKey: "parent.Field[0].InnerField.ErrorName", expected: "parent"},
		{name: "one_level", errorKey: "parent.Field.ErrorName", expected: "parent"},
		{name: "top_level", errorKey: "parent.ErrorName", expected: "parent"},
		{name: "empty", errorKey: "", expected: ""},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, tc.errorKey.TypeName())
		})
	}
}

func TestErrorKey_ErrorName(t *testing.T) {
	tcs := []struct {
		name     string
		errorKey ErrorKey
		expected string
	}{
		{name: "deep", errorKey: "parent.Field[0].InnerField.TheError", expected: "TheError"},
		{name: "one_level", errorKey: "parent.Field.TheError", expected: "TheError"},
		{name: "top_level", errorKey: "parent.TheError", expected: "TheError"},
		{name: "empty", errorKey: "", expected: ""},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, tc.errorKey.ErrorName())
		})
	}
}
