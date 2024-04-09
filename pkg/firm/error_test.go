package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorMap_Error(t *testing.T) {
	require := require.New(t)

	errorMap := ErrorMap{
		"A": TemplateError{Template: "field A message"},
		"B": TemplateError{Template: "field B message"},
	}
	require.Equal("A: value field A message, B: value field B message", errorMap.Error())
}

func TestErrorMap_MergeInto(t *testing.T) {
	require := require.New(t)

	src := ErrorMap{
		"A": TemplateError{Template: "a1"},
		"B": TemplateError{Template: "b1"},
	}
	dest := ErrorMap{
		"PATH.B": TemplateError{Template: "b2"},
		"PATH.C": TemplateError{Template: "c2"},
	}
	src.MergeInto("PATH", dest)
	require.Equal(ErrorMap{
		"PATH.A": TemplateError{Template: "a1"},
		"PATH.B": TemplateError{Template: "b1"},
		"PATH.C": TemplateError{Template: "c2"},
	}, dest)
}

func TestErrorMap_ToNil(t *testing.T) {
	tcs := []struct {
		name     string
		errorMap ErrorMap
		isNil    bool
	}{
		{name: "not_empty", errorMap: ErrorMap{"testy": TemplateError{}}},
		{name: "empty", errorMap: ErrorMap{}, isNil: true},
		{name: "nil", errorMap: nil, isNil: true},
	}
	for _, tc := range tcs {
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

var fullTemplateErrorKey = ErrorKey("pkger.Mover.Parent.MyField.TheError")
var errTemplateError = TemplateError{
	ErrorKey: fullTemplateErrorKey,
	Template: "has no {{ .Him }} and {{ .Her }} since it's of type: {{.RootTypeName}}",
	TemplateFields: map[string]string{
		"Him": "Jack",
		"Her": "Jill",
	},
}

func fullTemplate() TemplateError { return errTemplateError }

func TestErrorMap_Finish(t *testing.T) {
	tcs := []struct {
		name     string
		errorMap ErrorMap
		expected *TemplateError
	}{
		{name: "not_empty", errorMap: ErrorMap{fullTemplateErrorKey: TemplateError{}},
			expected: &TemplateError{ErrorKey: fullTemplateErrorKey}},
		{name: "empty", errorMap: ErrorMap{}, expected: nil},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			var expected ErrorMap
			if tc.expected != nil {
				expected = ErrorMap{fullTemplateErrorKey: *tc.expected}
			}
			require.Equal(expected, tc.errorMap.Finish())
		})
	}
}

func TestTemplateError_Error(t *testing.T) {
	tcs := []struct {
		name     string
		template string
		fields   map[string]string
		expected string
	}{
		{name: "everything", expected: "MyField has no Jack and Jill since it's of type: pkger.Mover"},
		{name: "empty_error_key", expected: "value has no Jack and Jill since it's of type: NoType"},
		{name: "missing_one", expected: "MyField has no <no value> and Jill since it's of type: pkger.Mover", fields: map[string]string{
			"Her": "Jill",
		}},
		{name: "bad_template", template: "{{ a }}", expected: "{{ a }} (bad format)"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			templateError := fullTemplate()
			if tc.name == "empty_error_key" {
				templateError.ErrorKey = ""
			}
			if tc.fields != nil {
				templateError.TemplateFields = tc.fields
			}
			if tc.template != "" {
				templateError.Template = tc.template
			}

			require.Equal(tc.expected, templateError.Error())
		})
	}
}

func TestErrorKey_RootTypeName(t *testing.T) {
	tcs := []struct {
		name     string
		errorKey ErrorKey
		expected string
	}{
		{name: "deep", errorKey: "firm.parent.Field.[0].InnerField.TheError", expected: "firm.parent"},
		{name: "one_level", errorKey: "firm.parent.Field.TheError", expected: "firm.parent"},
		{name: "one_level", errorKey: "firm.parent.Field.TheError", expected: "firm.parent"},
		{name: "top_level", errorKey: "firm.parent.TheError", expected: "firm.parent"},
		{name: "empty", errorKey: "", expected: ""},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, tc.errorKey.RootTypeName())
		})
	}
}

func TestErrorKey_ValueName(t *testing.T) {
	tcs := []struct {
		name     string
		errorKey ErrorKey
		expected string
	}{
		{name: "deep", errorKey: "firm.parent.Field.[0].InnerField.TheError", expected: "InnerField"},
		{name: "slice", errorKey: "firm.parent.Field.[0].TheError", expected: "[0]"},
		{name: "one_level", errorKey: "firm.parent.Field.TheError", expected: "Field"},
		{name: "top_level", errorKey: "firm.parent.TheError", expected: "firm.parent"},
		{name: "just_type", errorKey: "firm.parent", expected: ""},
		{name: "empty", errorKey: "", expected: ""},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, tc.errorKey.ValueName())
		})
	}
}

func TestErrorKey_ErrorName(t *testing.T) {
	tcs := []struct {
		name     string
		errorKey ErrorKey
		expected string
	}{
		{name: "deep", errorKey: "firm.parent.Field.[0].InnerField.TheError", expected: "TheError"},
		{name: "one_level", errorKey: "firm.parent.Field.TheError", expected: "TheError"},
		{name: "top_level", errorKey: "firm.parent.TheError", expected: "TheError"},
		{name: "empty", errorKey: "", expected: ""},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, tc.errorKey.ErrorName())
		})
	}
}

func TestRuleTypeError_TemplateError(t *testing.T) {
	require := require.New(t)
	require.Equal("value is not a string, got int",
		NewRuleTypeError("MyRule", reflect.TypeOf(0), "is not a string").TemplateError().Error())
}

func TestRuleTypeError_Error(t *testing.T) {
	require := require.New(t)
	require.Equal("MyRule: value is not a string, got int",
		NewRuleTypeError("MyRule", reflect.TypeOf(0), "is not a string").Error())
}
