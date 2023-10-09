package firm

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type registryParent struct {
	Primitive int
	Child     registryChild
}

type registryChild struct{}

type registryNotFoundTest struct{}

func TestRegistry_RegisterType(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	require.NoError(registry.RegisterType(NewDefinition[registryParent]().
		Validates(RuleMap{"Child": {}}).
		ValidatesTopLevel(presentRule{})))

	registryParentType := reflect.TypeOf(registryParent{})
	typeToValidator := map[reflect.Type]*ValueAny{
		registryParentType: {
			typ: registryParentType,
			rules: []Rule{presentRule{},
				StructAny{typ: reflect.TypeOf(registryParent{}), ruleMap: map[string]*[]Rule{"Child": {}}},
			}},
	}
	registryChildType := reflect.TypeOf(registryChild{})
	require.Equal(typeToValidator, registry.typeToValidator)
	require.Equal(map[reflect.Type][]*[]Rule{registryChildType: {{}}}, registry.unregisteredTypeRefs)

	require.NoError(registry.RegisterType(NewDefinition[registryChild]()))

	typeToValidator[registryParentType] = &ValueAny{typ: registryParentType, rules: []Rule{presentRule{},
		StructAny{typ: registryParentType, ruleMap: map[string]*[]Rule{
			"Child": {&ValueAny{typ: registryChildType, rules: []Rule{}}},
		}}}}
	typeToValidator[registryChildType] = &ValueAny{typ: registryChildType, rules: []Rule{}}
	require.Equal(typeToValidator, registry.typeToValidator)
	require.Equal(map[reflect.Type][]*[]Rule{}, registry.unregisteredTypeRefs)

	require.Equal(fmt.Errorf("RegisterType() with type firm.registryParent already exists"),
		registry.RegisterType(NewDefinition[registryParent]().ValidatesTopLevel(presentRule{})))
}

func notFoundError(data any) ErrorMap {
	value := reflect.ValueOf(data)
	errorMap := ErrorMap{}
	DefaultValidator.ValidateMerge(value, TypeName(value), errorMap)
	return errorMap.Finish()
}

// nolint:funlen // a bunch of test cases
func TestRegistry_ValidateAll(t *testing.T) {
	type testCase struct {
		name       string
		definition *Definition
		data       func() registryParent

		expectedKeySuffix string
		err               *TemplateError
	}
	tcs := []testCase{
		{
			name:              "top_level",
			definition:        NewDefinition[registryParent]().ValidatesTopLevel(presentRule{}),
			data:              func() registryParent { return registryParent{} },
			expectedKeySuffix: presentRuleKey,
			err:               presentRuleError(""),
		},
		{
			name: "field_Primitive",
			definition: NewDefinition[registryParent]().Validates(RuleMap{
				"Primitive": {presentRule{}},
			}),
			data:              func() registryParent { return registryParent{} },
			expectedKeySuffix: "Primitive.presentRule",
			err:               presentRuleError(""),
		},
		{
			name:              "not_found",
			definition:        NewDefinition[registryParent]().ValidatesTopLevel(presentRule{}),
			expectedKeySuffix: "NotFound",
		},
		{
			name: "not_found_field_Primitive",
			definition: NewDefinition[registryParent]().Validates(RuleMap{
				"Primitive": {presentRule{}},
			}),
			expectedKeySuffix: "NotFound",
		},
		{
			name:       "invalid",
			definition: NewDefinition[registryParent]().ValidatesTopLevel(presentRule{}),
			data:       nil,
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			registry := &Registry{}
			require.NoError(registry.RegisterType(tc.definition))

			if tc.name == "invalid" {
				var data any
				require.Equal(errInvalidValue, registry.ValidateAny(data))
				require.Equal(notFoundError(&data), registry.ValidateAny(&data))
				return
			}
			if strings.HasPrefix(tc.name, "not_found") {
				data := registryNotFoundTest{}
				require.Equal(notFoundError(data), registry.ValidateAny(data))
				require.Equal(notFoundError(&data), registry.ValidateAny(&data))

				notFoundTemplateError := &TemplateError{Template: "type, {{.RootTypeName}}, not found in Registry"}
				testValidateAllFull(t, true, registry, data, notFoundTemplateError, tc.expectedKeySuffix)
				testValidateAllFull(t, true, registry, &data, notFoundTemplateError, tc.expectedKeySuffix)
				return
			}
			data := tc.data()
			testValidateAll(t, registry, data, tc.err, tc.expectedKeySuffix)
			testValidateAll(t, registry, &data, tc.err, tc.expectedKeySuffix)
		})
	}
}

func TestRegistry_DefaultedValidator(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	require.NoError(registry.RegisterType(NewDefinition[registryParent]().ValidatesTopLevel(presentRule{})))
	expected, err := NewValueAny(reflect.TypeOf(registryParent{}), presentRule{})
	require.NoError(err)
	require.Equal(&expected, registry.DefaultedValidator(reflect.TypeOf(registryParent{})))

	notFoundType := reflect.TypeOf(nil)
	require.Equal(DefaultValidator, registry.DefaultedValidator(notFoundType))

	registry.DefaultValidator = Value[Any]{}
	require.Equal(registry.DefaultValidator, registry.DefaultedValidator(notFoundType))
}

func TestRegistry_Validator(t *testing.T) {
	registry := &Registry{}
	require.NoError(t, registry.RegisterType(NewDefinition[registryParent]().ValidatesTopLevel(presentRule{})))
	testParentValidator, err := NewValueAny(reflect.TypeOf(registryParent{}), presentRule{})
	require.NoError(t, err)

	tcs := []struct {
		name     string
		data     any
		expected Validator
	}{
		{name: "normal", data: registryParent{}, expected: &testParentValidator},
		{name: "pointer", data: &registryParent{}, expected: &testParentValidator},
		{name: "not_found", data: registryNotFoundTest{}},
		{name: "not_found_pointer", data: &registryNotFoundTest{}},
		{name: "nil_validator", data: Validator(nil)},
		{name: "pure_nil", data: nil},
		{name: "zero", data: 0},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expected, registry.Validator(reflect.TypeOf(tc.data)))
		})
	}
}
