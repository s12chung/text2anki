package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testPresence struct {
}

func (p testPresence) ValidateValue(value reflect.Value) ErrorMap {
	if !value.IsValid() || value.IsZero() {
		return ErrorMap{
			"testPresence": &TemplatedError{Template: "test"},
		}
	}
	return nil
}

func TestNotFoundRule_ValidateValue(t *testing.T) {
	require := require.New(t)
	value := reflect.ValueOf(1)
	errorMap := NotFoundRule{}.ValidateValue(value)
	require.Equal(ErrorMap{notFoundRuleErrorKey: notFoundRuleError(value)}, errorMap)
	require.Equal("type, int, not found in Registry", errorMap[notFoundRuleErrorKey].Error())
}

func TestMergeErrorMap(t *testing.T) {
	require := require.New(t)

	src := ErrorMap{
		"A": &TemplatedError{Template: "a1"},
		"B": &TemplatedError{Template: "b1"},
	}
	dest := ErrorMap{
		"PATH.B": &TemplatedError{Template: "b2"},
		"PATH.C": &TemplatedError{Template: "c2"},
	}
	MergeErrorMap("PATH", src, dest)
	require.Equal(ErrorMap{
		"PATH.A": &TemplatedError{Template: "a1"},
		"PATH.B": &TemplatedError{Template: "b1"},
		"PATH.C": &TemplatedError{Template: "c2"},
	}, dest)
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
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, tc.errorKey.ErrorName())
		})
	}
}
