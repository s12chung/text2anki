package firm

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type registryTestParent struct {
	Primitive int
	Child     registryTestChild
}

type registryTestChild struct{}

type registryNotFoundTest struct{}

func TestRegistry_RegisterType(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	require.NoError(registry.RegisterType(NewDefinition(registryTestParent{}).
		Validates(RuleMap{"Child": {}}).
		ValidatesTopLevel(presentRule{})))

	registryParentType := reflect.TypeOf(registryTestParent{})
	typeToValidator := map[reflect.Type]*ValueValidator{
		registryParentType: {
			typ: registryParentType,
			rules: []Rule{presentRule{},
				&StructValidator{typ: reflect.TypeOf(registryTestParent{}), ruleMap: map[string]*[]Rule{"Child": {}}},
			}},
	}
	registryChildType := reflect.TypeOf(registryTestChild{})
	require.Equal(typeToValidator, registry.typeToValidator)
	require.Equal(map[reflect.Type][]*[]Rule{registryChildType: {{}}}, registry.unregisteredTypeReferences)

	require.NoError(registry.RegisterType(NewDefinition(registryTestChild{})))

	typeToValidator[registryParentType] = &ValueValidator{typ: registryParentType, rules: []Rule{presentRule{},
		&StructValidator{typ: registryParentType, ruleMap: map[string]*[]Rule{
			"Child": {&ValueValidator{typ: registryChildType, rules: []Rule{&StructValidator{typ: registryChildType, ruleMap: map[string]*[]Rule{}}}}},
		}}}}
	typeToValidator[registryChildType] = &ValueValidator{typ: registryChildType,
		rules: []Rule{&StructValidator{typ: registryChildType, ruleMap: map[string]*[]Rule{}}},
	}
	require.Equal(typeToValidator, registry.typeToValidator)
	require.Equal(map[reflect.Type][]*[]Rule{}, registry.unregisteredTypeReferences)

	require.Equal(fmt.Errorf("RegisterType() with type firm.registryTestParent already exists"),
		registry.RegisterType(NewDefinition(registryTestParent{}).ValidatesTopLevel(presentRule{})))
}

func notFoundError(data any) ErrorMap {
	value := reflect.ValueOf(data)
	errorMap := ErrorMap{}
	DefaultValidator.ValidateMerge(value, TypeName(value), errorMap)
	return errorMap.Finish()
}

// nolint:funlen // a bunch of test cases
func TestRegistry_Validate(t *testing.T) {
	type testCase struct {
		name       string
		definition *Definition
		data       func() registryTestParent

		expectedKeySuffix string
		err               *TemplateError
	}
	tcs := []testCase{
		{
			name:              "top_level",
			definition:        NewDefinition(registryTestParent{}).ValidatesTopLevel(presentRule{}),
			data:              func() registryTestParent { return registryTestParent{} },
			expectedKeySuffix: presentRuleKey,
			err:               presentRuleError(),
		},
		{
			name: "field_Primitive",
			definition: NewDefinition(registryTestParent{}).
				Validates(RuleMap{
					"Primitive": {presentRule{}},
				}),
			data:              func() registryTestParent { return registryTestParent{} },
			expectedKeySuffix: "Primitive.presentRule",
			err:               presentRuleError(),
		},
		{
			name:              "not_found",
			definition:        NewDefinition(registryTestParent{}).ValidatesTopLevel(presentRule{}),
			expectedKeySuffix: "NotFound",
		},
		{
			name: "not_found_field_Primitive",
			definition: NewDefinition(registryTestParent{}).
				Validates(RuleMap{
					"Primitive": {presentRule{}},
				}),
			expectedKeySuffix: "NotFound",
		},
		{
			name:       "invalid",
			definition: NewDefinition(registryTestParent{}).ValidatesTopLevel(presentRule{}),
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
				require.Equal(errInvalidValue, registry.Validate(data))
				require.Equal(notFoundError(&data), registry.Validate(&data))
				return
			}
			if strings.HasPrefix(tc.name, "not_found") {
				data := registryNotFoundTest{}
				require.Equal(notFoundError(data), registry.Validate(data))
				require.Equal(notFoundError(&data), registry.Validate(&data))

				notFoundTemplateError := &TemplateError{Template: "type, {{.RootTypeName}}, not found in Registry"}
				testValidatesFull(t, true, registry, data, notFoundTemplateError, tc.expectedKeySuffix)
				testValidatesFull(t, true, registry, &data, notFoundTemplateError, tc.expectedKeySuffix)
				return
			}
			data := tc.data()
			testValidates(t, registry, data, tc.err, tc.expectedKeySuffix)
			testValidates(t, registry, &data, tc.err, tc.expectedKeySuffix)
		})
	}
}

func TestRegistry_DefaultedValidator(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	require.NoError(registry.RegisterType(NewDefinition(registryTestParent{}).ValidatesTopLevel(presentRule{})))

	registryParentType := reflect.TypeOf(registryTestParent{})
	structValidator, err := newStructValidator(registryParentType, nil)
	require.NoError(err)
	expected, err := newValueValidator(registryParentType, presentRule{}, &structValidator)
	require.NoError(err)
	require.Equal(&expected, registry.DefaultedValidator(registryParentType))

	notFoundType := reflect.TypeOf(nil)
	require.Equal(DefaultValidator, registry.DefaultedValidator(notFoundType))

	registry.DefaultValidator = ValueValidator{}
	require.Equal(Validator(ValueValidator{}), registry.DefaultedValidator(notFoundType))
}

func TestRegistry_Validator(t *testing.T) {
	registry := &Registry{}
	require.NoError(t, registry.RegisterType(NewDefinition(registryTestParent{}).ValidatesTopLevel(presentRule{})))

	registryParentType := reflect.TypeOf(registryTestParent{})
	structValidator, err := newStructValidator(registryParentType, nil)
	require.NoError(t, err)
	testParentValidator, err := newValueValidator(registryParentType, presentRule{}, &structValidator)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		data     any
		expected Validator
	}{
		{name: "normal", data: registryTestParent{}, expected: &testParentValidator},
		{name: "pointer", data: &registryTestParent{}, expected: &testParentValidator},
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
