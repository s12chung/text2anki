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
	errorKeys []ErrorKey
}

var structValidatorTestCases = []structValidatorTestCase{
	//
	// Full
	//
	{name: "Full", errorKeys: nil, f: fullParent},

	//
	// Primitive
	//
	{name: "Primitive___zero", errorKeys: []ErrorKey{"Primitive"}, f: func() parent {
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
	{name: "Basic___child_validates_zero", errorKeys: []ErrorKey{"Basic.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Basic.Validates = ""
		return changeParent
	}},
	{name: "Basic___child_empty", errorKeys: []ErrorKey{"Basic", "Basic.Validates"}, f: func() parent {
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
	{name: "Pt___child_validates_zero", errorKeys: []ErrorKey{"Pt.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Pt.Validates = ""
		return changeParent
	}},
	{name: "Pt___child_empty", errorKeys: []ErrorKey{"Pt.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Pt = &child{}
		return changeParent
	}},
	{name: "Pt___nil", errorKeys: []ErrorKey{"Pt"}, f: func() parent {
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
	{name: "Any___nil", errorKeys: []ErrorKey{"Any"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Any = nil
		return changeParent
	}},

	//
	// Array
	//
	{name: "Array___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		for _, v := range changeParent.Array {
			v.NoValidates = ""
		}
		return changeParent
	}},
	{name: "Array___child_validates_zero", errorKeys: []ErrorKey{"Array[0].Validates", "Array[1].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Array[0].Validates = ""
		changeParent.Array[1].Validates = ""
		return changeParent
	}},
	{name: "Array___child_validates_one_zero", errorKeys: []ErrorKey{"Array[0].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Array[0].Validates = ""
		return changeParent
	}},
	{name: "Array___empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.Array = []child{}
		return changeParent
	}},
	{name: "Array___nil", errorKeys: []ErrorKey{"Array"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Array = nil
		return changeParent
	}},

	//
	// ArrayPt
	//
	{name: "ArrayPt___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		for _, v := range changeParent.ArrayPt {
			v.NoValidates = ""
		}
		return changeParent
	}},
	{name: "ArrayPt___child_validates_zero", errorKeys: []ErrorKey{"ArrayPt[0].Validates", "ArrayPt[1].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPt[0].Validates = ""
		changeParent.ArrayPt[1].Validates = ""
		return changeParent
	}},
	{name: "ArrayPt___child_validates_one_zero", errorKeys: []ErrorKey{"ArrayPt[0].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPt[0].Validates = ""
		return changeParent
	}},
	{name: "ArrayPt___empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPt = []*child{}
		return changeParent
	}},
	{name: "ArrayPt___nil", errorKeys: []ErrorKey{"ArrayPt"}, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPt = nil
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
	{name: "BasicEmptyValidates___child_validates_zero", errorKeys: []ErrorKey{"BasicEmptyValidates.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.BasicEmptyValidates.Validates = ""
		return changeParent
	}},
	{name: "BasicEmptyValidates____child_empty", errorKeys: []ErrorKey{"BasicEmptyValidates.Validates"}, f: func() parent {
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
	{name: "PtEmptyValidates___child_validates_zero", errorKeys: []ErrorKey{"PtEmptyValidates.Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.PtEmptyValidates.Validates = ""
		return changeParent
	}},
	{name: "PtEmptyValidates___child_empty", errorKeys: []ErrorKey{"PtEmptyValidates.Validates"}, f: func() parent {
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
	// ArrayValidates
	//
	{name: "ArrayValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		for _, v := range changeParent.ArrayValidates {
			v.NoValidates = ""
		}
		return changeParent
	}},
	{name: "ArrayValidates___child_validates_zero",
		errorKeys: []ErrorKey{"ArrayValidates[0].Validates", "ArrayValidates[1].Validates"}, f: func() parent {
			changeParent := fullParent()
			changeParent.ArrayValidates[0].Validates = ""
			changeParent.ArrayValidates[1].Validates = ""
			return changeParent
		}},
	{name: "ArrayValidates___child_validates_one_zero", errorKeys: []ErrorKey{"ArrayValidates[0].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayValidates[0].Validates = ""
		return changeParent
	}},
	{name: "ArrayValidates___empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayValidates = []child{}
		return changeParent
	}},
	{name: "ArrayValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayValidates = nil
		return changeParent
	}},

	//
	// ArrayPtValidates
	//
	{name: "ArrayPtValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		for _, v := range changeParent.ArrayPtValidates {
			v.NoValidates = ""
		}
		return changeParent
	}},
	{name: "ArrayPtValidates___child_validates_zero",
		errorKeys: []ErrorKey{"ArrayPtValidates[0].Validates", "ArrayPtValidates[1].Validates"}, f: func() parent {
			changeParent := fullParent()
			changeParent.ArrayPtValidates[0].Validates = ""
			changeParent.ArrayPtValidates[1].Validates = ""
			return changeParent
		}},
	{name: "ArrayPtValidates___child_validates_one_zero", errorKeys: []ErrorKey{"ArrayPtValidates[0].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtValidates[0].Validates = ""
		return changeParent
	}},
	{name: "ArrayPtValidates___empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtValidates = []*child{}
		return changeParent
	}},
	{name: "ArrayPtValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtValidates = nil
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

	//
	// ArrayNoValidates
	//
	{name: "ArrayNoValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		for _, v := range changeParent.ArrayNoValidates {
			v.NoValidates = ""
		}
		return changeParent
	}},
	{name: "ArrayNoValidates___child_validates_zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayNoValidates[0].Validates = ""
		changeParent.ArrayNoValidates[1].Validates = ""
		return changeParent
	}},
	{name: "ArrayNoValidates___child_validates_one_zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayNoValidates[0].Validates = ""
		return changeParent
	}},
	{name: "ArrayNoValidates___empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayNoValidates = []child{}
		return changeParent
	}},
	{name: "ArrayNoValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayNoValidates = nil
		return changeParent
	}},

	//
	// ArrayPtNoValidates
	//
	{name: "ArrayPtNoValidates___child_validates_ok", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		for _, v := range changeParent.ArrayPtNoValidates {
			v.NoValidates = ""
		}
		return changeParent
	}},
	{name: "ArrayPtNoValidates___child_validates_zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtNoValidates[0].Validates = ""
		changeParent.ArrayPtNoValidates[1].Validates = ""
		return changeParent
	}},
	{name: "ArrayPtNoValidates___child_validates_one_zero", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtNoValidates[0].Validates = ""
		return changeParent
	}},
	{name: "ArrayPtNoValidates___empty", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtNoValidates = []*child{}
		return changeParent
	}},
	{name: "ArrayPtNoValidates___nil", errorKeys: nil, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPtNoValidates = nil
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
				key := joinKeys(typeNameKey(reflect.ValueOf(rawData)), errorKey, "testPresence")
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
				expected.errorMap[joinKeys(typeNameKey(reflect.ValueOf(tc.data)), "testPresence")] = tc.err
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
