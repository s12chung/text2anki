package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

type structValidatorTestCase struct {
	name      string
	f         func() parent
	errorKeys []string
}

var structValidatorTestCases = []structValidatorTestCase{
	//
	// Full
	//
	{name: "Full", errorKeys: nil, f: fullParent},

	//
	// Primitive
	//
	{name: "Primitive___zero", errorKeys: []string{"Primitive"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Primitive = 0
		return changeParent
	}},

	//
	// Basic
	//
	{name: "Basic___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.Basic.NoValidates = ""
		return changeParent
	}},
	{name: "Basic___child_validates_zero", errorKeys: []string{"Basic.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Basic.Validates = ""
		return changeParent
	}},
	{name: "Basic___child_empty", errorKeys: []string{"Basic", "Basic.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Basic = child{}
		return changeParent
	}},

	//
	// Pt
	//
	{name: "Pt___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.Pt.NoValidates = ""
		return changeParent
	}},
	{name: "Pt___child_validates_zero", errorKeys: []string{"Pt.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Pt.Validates = ""
		return changeParent
	}},
	{name: "Pt___child_empty", errorKeys: []string{"Pt.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Pt = &child{}
		return changeParent
	}},
	{name: "Pt___nil", errorKeys: []string{"Pt"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Pt = nil
		return changeParent
	}},

	//
	// Any
	//
	{name: "Any___child_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.Any = child{}
		return changeParent
	}},
	{name: "Any___child_pointer_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.Any = &child{}
		return changeParent
	}},
	{name: "Any___nil", errorKeys: []string{"Any"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Any = nil
		return changeParent
	}},

	//
	// PrimitiveEmptyValidates
	//
	{name: "PrimitiveEmptyValidates___zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PrimitiveEmptyValidates = 0
		return changeParent
	}},

	//
	// BasicEmptyValidates
	//
	{name: "BasicEmptyValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicEmptyValidates.NoValidates = ""
		return changeParent
	}},
	{name: "BasicEmptyValidates___child_validates_zero", errorKeys: []string{"BasicEmptyValidates.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicEmptyValidates.Validates = ""
		return changeParent
	}},
	{name: "BasicEmptyValidates____child_empty", errorKeys: []string{"BasicEmptyValidates.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicEmptyValidates = child{}
		return changeParent
	}},

	//
	// PtEmptyValidates
	//
	{name: "PtEmptyValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PtEmptyValidates.NoValidates = ""
		return changeParent
	}},
	{name: "PtEmptyValidates___child_validates_zero", errorKeys: []string{"PtEmptyValidates.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.PtEmptyValidates.Validates = ""
		return changeParent
	}},
	{name: "PtEmptyValidates___child_empty", errorKeys: []string{"PtEmptyValidates.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.PtEmptyValidates = &child{}
		return changeParent
	}},
	{name: "PtEmptyValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PtEmptyValidates = nil
		return changeParent
	}},

	//
	// AnyEmptyValidates
	//
	{name: "AnyEmptyValidates___child_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.AnyEmptyValidates = child{}
		return changeParent
	}},
	{name: "AnyEmptyValidates___child_pointer_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.AnyEmptyValidates = &child{}
		return changeParent
	}},
	{name: "AnyEmptyValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.AnyEmptyValidates = nil
		return changeParent
	}},

	//
	// PrimitiveNoValidates
	//
	{name: "PrimitiveNoValidates___zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PrimitiveNoValidates = 0
		return changeParent
	}},

	//
	// BasicNoValidates
	//
	{name: "BasicNoValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicNoValidates.NoValidates = ""
		return changeParent
	}},
	{name: "BasicNoValidates___child_validates_zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicNoValidates.Validates = ""
		return changeParent
	}},
	{name: "BasicNoValidates____child_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicNoValidates = child{}
		return changeParent
	}},

	//
	// PtNoValidates
	//
	{name: "PtNoValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PtNoValidates.NoValidates = ""
		return changeParent
	}},
	{name: "PtNoValidates___child_validates_zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PtNoValidates.Validates = ""
		return changeParent
	}},
	{name: "PtNoValidates___child_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PtNoValidates = &child{}
		return changeParent
	}},
	{name: "PtNoValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.PtNoValidates = nil
		return changeParent
	}},

	//
	// AnyNoValidates
	//
	{name: "AnyNoValidates___child_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.AnyNoValidates = child{}
		return changeParent
	}},
	{name: "AnyNoValidates___child_pointer_empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.AnyNoValidates = &child{}
		return changeParent
	}},
	{name: "AnyNoValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.AnyNoValidates = nil
		return changeParent
	}},
}

func TestStructValidator_Validate(t *testing.T) {
	validator := testRegistry.DefinitionForType(reflect.TypeOf(parent{})).Validator(testRegistry)

	t.Run("not_struct", func(t *testing.T) {
		require := require.New(t)

		data := 1
		dataValue := reflect.ValueOf(data)
		expected := MapResult{errorMap: ErrorMap{
			joinKeys("int", structValidatorErrorKey): structValidatorError(dataValue),
		}}
		require.Equal(expected, validator.Validate(data))
		require.Equal(expected.errorMap, validator.ValidateValue(dataValue))
		errorMap := ErrorMap{
			"Existing": nil,
		}
		expectedErrorMap := maps.Clone(errorMap)
		expectedErrorMap[structValidatorErrorKey] = structValidatorError(dataValue)
		validator.ValidateMerge(dataValue, "KEY.ME", errorMap)
	})

	for _, tc := range structValidatorTestCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			rawData := tc.f()
			expected := MapResult{errorMap: ErrorMap{}}
			for _, errorKey := range tc.errorKeys {
				key := joinKeys(typeName(reflect.ValueOf(rawData)), errorKey, "testPresence")
				expected.errorMap[key] = &TemplatedError{Template: "test"}
			}

			for _, data := range []any{rawData, &rawData} {
				dataValue := reflect.ValueOf(data)
				require.Equal(expected, validator.Validate(data))
				require.Equal(expected.errorMap, validator.ValidateValue(dataValue))

				errorMap := ErrorMap{
					"Existing": nil,
				}
				expectedErrorMap := maps.Clone(errorMap)
				for _, errorKey := range tc.errorKeys {
					key := joinKeys("KEY.ME", errorKey, "testPresence")
					expectedErrorMap[key] = &TemplatedError{Template: "test"}
				}
				validator.ValidateMerge(dataValue, "KEY.ME", errorMap)
				require.Equal(expectedErrorMap, errorMap)
			}
		})
	}
}

func TestValueValidator_Validate(t *testing.T) {
	type testCase struct {
		name string
		data any

		err *TemplatedError
	}
	tcs := []testCase{
		{
			name: "not_zero",
			data: 1,
		},
		{
			name: "zero",
			data: 0,
			err:  &TemplatedError{Template: "test"},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			validator := NewValueValidator(testPresence{})

			expected := MapResult{errorMap: ErrorMap{}}
			if tc.err != nil {
				expected.errorMap[joinKeys(typeName(reflect.ValueOf(tc.data)), "testPresence")] = tc.err
			}
			require.Equal(expected, validator.Validate(tc.data))
			require.Equal(expected.errorMap, validator.ValidateValue(reflect.ValueOf(tc.data)))

			errorMap := ErrorMap{
				"Existing": nil,
			}
			expectedErrorMap := maps.Clone(errorMap)
			if tc.err != nil {
				expectedErrorMap[joinKeys("KEY.ME", "testPresence")] = tc.err
			}
			validator.ValidateMerge(reflect.ValueOf(tc.data), "KEY.ME", errorMap)
			require.Equal(expectedErrorMap, errorMap)
		})
	}
}
