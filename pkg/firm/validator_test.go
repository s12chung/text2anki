package firm

import (
	"maps"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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
	{name: "Array___child_validates_zero", errorKeys: []ErrorKey{"Array.[0].Validates", "Array.[1].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.Array[0].Validates = ""
		changeParent.Array[1].Validates = ""
		return changeParent
	}},
	{name: "Array___child_validates_one_zero", errorKeys: []ErrorKey{"Array.[0].Validates"}, f: func() parent {
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
	{name: "ArrayPt___child_validates_zero", errorKeys: []ErrorKey{"ArrayPt.[0].Validates", "ArrayPt.[1].Validates"}, f: func() parent {
		changeParent := fullParent()
		changeParent.ArrayPt[0].Validates = ""
		changeParent.ArrayPt[1].Validates = ""
		return changeParent
	}},
	{name: "ArrayPt___child_validates_one_zero", errorKeys: []ErrorKey{"ArrayPt.[0].Validates"}, f: func() parent {
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
		errorKeys: []ErrorKey{"ArrayValidates.[0].Validates", "ArrayValidates.[1].Validates"}, f: func() parent {
			changeParent := fullParent()
			changeParent.ArrayValidates[0].Validates = ""
			changeParent.ArrayValidates[1].Validates = ""
			return changeParent
		}},
	{name: "ArrayValidates___child_validates_one_zero", errorKeys: []ErrorKey{"ArrayValidates.[0].Validates"}, f: func() parent {
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
		errorKeys: []ErrorKey{"ArrayPtValidates.[0].Validates", "ArrayPtValidates.[1].Validates"}, f: func() parent {
			changeParent := fullParent()
			changeParent.ArrayPtValidates[0].Validates = ""
			changeParent.ArrayPtValidates[1].Validates = ""
			return changeParent
		}},
	{name: "ArrayPtValidates___child_validates_one_zero", errorKeys: []ErrorKey{"ArrayPtValidates.[0].Validates"}, f: func() parent {
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
	validator := testRegistry.Validator(reflect.TypeOf(parent{}))

	tcs := []struct {
		name   string
		data   any
		result ErrorMap
	}{
		{name: "not_struct", data: 1, result: validateTypeErrorResult(validator, 1)},
		{name: "invalid", data: nil, result: errInvalidValue},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { require.Equal(t, tc.result, validator.Validate(tc.data)) })
	}

	for _, tc := range structValidatorTestCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rawData := tc.f()
			errKeys := make([]ErrorKey, len(tc.errorKeys))
			for i, key := range tc.errorKeys {
				errKeys[i] = joinKeys(key, presentRuleKey)
			}
			testValidates(t, validator, rawData, errTest, errKeys...)
			testValidates(t, validator, &rawData, errTest, errKeys...)
		})
	}
}

func TestStructValidator_ValidateType(t *testing.T) {
	parentType := reflect.TypeOf(parent{})
	validator := testRegistry.Validator(parentType)
	badCondition := "is not matching Struct of type firm.parent"

	tcs := []struct {
		name         string
		typ          reflect.Type
		badCondition string
	}{
		{name: "matching struct", typ: parentType},
		{name: "matching struct pointer", typ: reflect.TypeOf(&parent{})},
		{name: "other struct", typ: reflect.TypeOf(child{}), badCondition: badCondition},
		{name: "not struct", typ: reflect.TypeOf(1), badCondition: badCondition},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			var err *RuleTypeError
			if tc.badCondition != "" {
				err = NewRuleTypeError(tc.typ, tc.badCondition)
			}
			require.Equal(err, validator.ValidateType(tc.typ))
		})
	}
}

type sliceValidatorElement struct {
	Int  int
	UInt uint
}

type sliceValidatorTestCase struct {
	name      string
	f         func() []sliceValidatorElement
	errorKeys []ErrorKey
}

var sliceValidatorTestCases = []sliceValidatorTestCase{
	{name: "Full", errorKeys: nil, f: func() []sliceValidatorElement {
		return []sliceValidatorElement{{1, 1}, {2, 2}, {3, 3}}
	}},
	{name: "Empty", errorKeys: nil, f: func() []sliceValidatorElement {
		return []sliceValidatorElement{}
	}},
	{name: "Nil", errorKeys: nil, f: func() []sliceValidatorElement {
		return nil
	}},
	{name: "Element_Not_Full", errorKeys: nil, f: func() []sliceValidatorElement {
		return []sliceValidatorElement{{Int: 1}, {Int: 2}}
	}},
	{name: "Element_Invalid", errorKeys: []ErrorKey{"[0].Int", "[1].Int"}, f: func() []sliceValidatorElement {
		return []sliceValidatorElement{{UInt: 1}, {UInt: 2}}
	}},
	{name: "Element_Empty", errorKeys: []ErrorKey{"[0]", "[1]", "[0].Int", "[1].Int"}, f: func() []sliceValidatorElement {
		return []sliceValidatorElement{{}, {}}
	}},
}

var sliceValidator = NewSliceValidator(reflect.TypeOf([]sliceValidatorElement{}),
	presentRule{}, NewStructValidator(reflect.TypeOf([]sliceValidatorElement{}), RuleMap{"Int": {presentRule{}}}))

func TestSliceValidator_Validate(t *testing.T) {
	validator := sliceValidator

	tcs := []struct {
		name   string
		data   any
		result ErrorMap
	}{
		{name: "not_slice", data: 1, result: validateTypeErrorResult(validator, 1)},
		{name: "invalid", data: nil, result: errInvalidValue},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { require.Equal(t, tc.result, validator.Validate(tc.data)) })
	}

	for _, tc := range sliceValidatorTestCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rawData := tc.f()
			errKeys := make([]ErrorKey, len(tc.errorKeys))
			for i, key := range tc.errorKeys {
				errKeys[i] = joinKeys(key, presentRuleKey)
			}
			testValidates(t, validator, rawData, errTest, errKeys...)
			testValidates(t, validator, &rawData, errTest, errKeys...)
		})
	}
}

func TestSliceValidator_ValidateType(t *testing.T) {
	sliceType := reflect.TypeOf([]sliceValidatorElement{})
	validator := sliceValidator
	badCondition := "is not matching Slice or Array of type []firm.sliceValidatorElement"

	tcs := []struct {
		name         string
		typ          reflect.Type
		badCondition string
	}{
		{name: "matching slice", typ: sliceType},
		{name: "matching slice pointer", typ: reflect.TypeOf(&[]sliceValidatorElement{})},
		{name: "other slice", typ: reflect.TypeOf([]int{}), badCondition: badCondition},
		{name: "not slice", typ: reflect.TypeOf(1), badCondition: badCondition},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			var err *RuleTypeError
			if tc.badCondition != "" {
				err = NewRuleTypeError(tc.typ, tc.badCondition)
			}
			require.Equal(err, validator.ValidateType(tc.typ))
		})
	}
}

type onlyKindRule struct{ kind reflect.Kind }

func (o onlyKindRule) ValidateValue(_ reflect.Value) ErrorMap { return nil }
func (o onlyKindRule) ValidateType(typ reflect.Type) *RuleTypeError {
	if typ.Kind() != o.kind {
		return NewRuleTypeError(typ, "is not "+o.kind.String())
	}
	return nil
}

func TestValueValidator_Validate(t *testing.T) {
	validator := NewValueValidator(presentRule{})
	onlyKindValidator := NewValueValidator(onlyKindRule{kind: reflect.String})

	edgeTcs := []struct {
		name      string
		validator Validator
		data      any
		result    ErrorMap
	}{
		{name: "not_valid_type", validator: onlyKindValidator, data: 1, result: validateTypeErrorResult(onlyKindValidator, 1)},
		{name: "invalid", validator: validator, data: nil, result: errInvalidValue},
	}
	for _, tc := range edgeTcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { require.Equal(t, tc.result, tc.validator.Validate(tc.data)) })
	}

	type testCase struct {
		name string
		data any
		err  *TemplatedError
	}
	tcs := []testCase{
		{name: "not_zero", data: 1},
		{name: "zero", data: 0, err: errTest},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { testValidates(t, validator, tc.data, tc.err, presentRuleKey) })
	}
}

func TestValueValidator_ValidateType(t *testing.T) {
	intType := reflect.TypeOf(1)
	badCondition := "is not int"

	i := 0
	tcs := []struct {
		name         string
		typ          reflect.Type
		extraRule    Rule
		badCondition string
	}{
		{name: "matching int", typ: intType},
		{name: "matching int pointer", typ: reflect.TypeOf(&i), badCondition: badCondition},
		{name: "not int", typ: reflect.TypeOf([]int{}), badCondition: badCondition},
		{name: "extra rule", typ: intType, extraRule: onlyKindRule{kind: reflect.String}, badCondition: "is not string"},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			rules := []Rule{presentRule{}, onlyKindRule{kind: reflect.Int}, presentRule{}}
			if tc.extraRule != nil {
				rules = append(rules, tc.extraRule)
			}
			validator := NewValueValidator(rules...)

			var err *RuleTypeError
			if tc.badCondition != "" {
				err = NewRuleTypeError(tc.typ, tc.badCondition)
			}
			require.Equal(err, validator.ValidateType(tc.typ))
		})
	}
}

func testValidates(t *testing.T, validator Validator, data any, err *TemplatedError, keySuffixes ...ErrorKey) {
	testValidatesFull(t, false, validator, data, err, keySuffixes...)
}

func testValidatesFull(t *testing.T, skipValidate bool, validator Validator, data any, err *TemplatedError, keySuffixes ...ErrorKey) {
	require := require.New(t)

	var validateValueExpected ErrorMap
	if err != nil && len(keySuffixes) > 0 {
		validateValueExpected = ErrorMap{}
		for _, key := range keySuffixes {
			validateValueExpected[key] = err
		}
	}
	validateExpected := ErrorMap{}
	validateValueExpected.MergeInto(typeNameKey(reflect.ValueOf(data)), validateExpected)
	validateExpected = validateExpected.ToNil()

	if !skipValidate {
		require.Equal(validateExpected, validator.Validate(data))
	}
	require.Equal(validateValueExpected, validator.ValidateValue(reflect.ValueOf(data)))

	errorMap := ErrorMap{"Existing": nil}
	expectedErrorMap := maps.Clone(errorMap)
	if err != nil {
		for _, key := range keySuffixes {
			expectedErrorMap[joinKeys("KEY.ME", key)] = err
		}
	}
	validator.ValidateMerge(reflect.ValueOf(data), "KEY.ME", errorMap)
	require.Equal(expectedErrorMap, errorMap)
}
