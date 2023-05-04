package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func TestValueValidator_Validate(t *testing.T) {
	type testCase struct {
		name string
		data any

		expectedKeySuffix string
		err               *TemplatedError
	}
	tcs := []testCase{
		{
			name: "not_zero",
			data: 1,
		},
		{
			name:              "zero",
			data:              0,
			expectedKeySuffix: "testPresence",
			err:               &TemplatedError{Template: "test"},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			validator := NewValueValidator(testPresence{})

			expected := MapResult{errorMap: ErrorMap{}}
			if tc.err != nil {
				expected.errorMap = ErrorMap{
					joinKeys(typeName(reflect.ValueOf(tc.data)), tc.expectedKeySuffix): tc.err,
				}
			}
			require.Equal(expected, validator.Validate(tc.data))
			require.Equal(expected.errorMap, validator.ValidateValue(reflect.ValueOf(tc.data)))

			errorMap := ErrorMap{
				"Existing": nil,
			}
			expectedErrorMap := maps.Clone(errorMap)
			if tc.err != nil {
				expectedErrorMap[joinKeys("KEY.ME", tc.expectedKeySuffix)] = tc.err
			}
			validator.ValidateMerge(reflect.ValueOf(tc.data), "KEY.ME", errorMap)
			require.Equal(expectedErrorMap, errorMap)
		})
	}
}
