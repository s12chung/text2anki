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
		Key:      "KEY.ME",
		Template: "the error we go with {{ .Him }} and {{ .Her }}. yay",
		TemplateFields: map[string]string{
			"Him": "Jack",
			"Her": "Jill",
		},
	}
	require.Equal("invalid KEY.ME: the error we go with Jack and Jill. yay", err.Error())

	err.TemplateFields = map[string]string{
		"Her": "Jill",
	}
	require.Equal("invalid KEY.ME: the error we go with <no value> and Jill. yay", err.Error())

	err.Template = "{{ a }}"
	require.Equal("invalid KEY.ME: {{ a }} (bad format)", err.Error())
}
