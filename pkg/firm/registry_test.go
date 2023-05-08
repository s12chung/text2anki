package firm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

type registryTest struct {
	Primitive int
}

type registryNotFoundTest struct{}

func TestRegistry_RegisterType(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	definition := NewTypedDefinition(registryTest{})

	registry.RegisterType(definition)
	require.Equal(map[reflect.Type]StructuredDefinition{
		definition.typ: definition,
	}, registry.typeToDefinition)

	require.Panics(func() {
		registry.RegisterType(NewTypedDefinition(registryTest{}).ValidatesTopLevel(testPresence{}))
	})
}

func TestRegistry_DefinitionForType(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	definition := NewTypedDefinition(registryTest{})
	require.Nil(registry.DefinitionForType(definition.typ))

	registry.RegisterType(definition)
	require.Equal(definition, registry.DefinitionForType(definition.typ))
	require.Nil(registry.DefinitionForType(nil))
	require.Nil(registry.DefinitionForType(reflect.TypeOf(1)))
}

// nolint:funlen // a bunch of test cases
func TestRegistry_Validate(t *testing.T) {
	type testCase struct {
		name       string
		definition *TypedDefinition
		data       any

		expectedKeySuffix ErrorKey
		err               *TemplatedError
	}
	tcs := []testCase{
		{
			name:              "top_level",
			definition:        NewTypedDefinition(registryTest{}).ValidatesTopLevel(testPresence{}),
			data:              registryTest{},
			expectedKeySuffix: "testPresence",
			err:               &TemplatedError{Template: "test"},
		},
		{
			name: "field_Primitive",
			definition: NewTypedDefinition(registryTest{}).
				Validates(RuleMap{
					"Primitive": {testPresence{}},
				}),
			data:              registryTest{},
			expectedKeySuffix: "Primitive.testPresence",
			err:               &TemplatedError{Template: "test"},
		},
		{
			name:              "not_found",
			definition:        NewTypedDefinition(registryTest{}).ValidatesTopLevel(testPresence{}),
			data:              registryNotFoundTest{},
			expectedKeySuffix: notFoundRuleErrorKey,
			err:               notFoundRuleError(reflect.ValueOf(registryNotFoundTest{})),
		},
		{
			name: "not_found_field_Primitive",
			definition: NewTypedDefinition(registryTest{}).
				Validates(RuleMap{
					"Primitive": {testPresence{}},
				}),
			data:              registryNotFoundTest{},
			expectedKeySuffix: notFoundRuleErrorKey,
			err:               notFoundRuleError(reflect.ValueOf(registryNotFoundTest{})),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			registry := &Registry{}
			registry.RegisterType(tc.definition)

			expected := MapResult{errorMap: ErrorMap{}}
			if tc.err != nil {
				expected.errorMap[joinKeys(typeNameKey(reflect.ValueOf(tc.data)), tc.expectedKeySuffix)] = tc.err
			}
			require.Equal(expected, registry.Validate(tc.data))
			require.Equal(expected.errorMap, registry.ValidateValue(reflect.ValueOf(tc.data)))

			errorMap := ErrorMap{
				"Existing": nil,
			}
			expectedErrorMap := maps.Clone(errorMap)
			if tc.err != nil {
				expectedErrorMap[joinKeys("KEY.ME", tc.expectedKeySuffix)] = tc.err
			}
			registry.ValidateMerge(reflect.ValueOf(tc.data), "KEY.ME", errorMap)
			require.Equal(expectedErrorMap, errorMap)
		})
	}
}

func TestRegistry_DefaultedValidator(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	definition := NewTypedDefinition(registryTest{})
	registry.RegisterType(definition)

	require.Equal(definition.Validator(registry), registry.DefaultedValidator(reflect.ValueOf(registryTest{})))

	notFoundValue := reflect.ValueOf(nil)
	require.Equal(DefaultValidator, registry.DefaultedValidator(notFoundValue))

	registry.DefaultValidator = ValueValidator{}
	require.Equal(Validator(ValueValidator{}), registry.DefaultedValidator(notFoundValue))
}

func TestRegistry_Validator(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	definition := NewTypedDefinition(registryTest{})
	registry.RegisterType(definition)

	require.Equal(definition.Validator(registry), registry.Validator(reflect.ValueOf(registryTest{})))
	require.Equal(definition.Validator(registry), registry.Validator(reflect.ValueOf(&registryTest{})))
	require.Nil(registry.Validator(reflect.ValueOf(registryNotFoundTest{})))
	require.Nil(registry.Validator(reflect.ValueOf(&registryNotFoundTest{})))
	require.Nil(registry.Validator(reflect.ValueOf(Validator(nil))))
	require.Nil(registry.Validator(reflect.ValueOf(nil)))
	require.Nil(registry.Validator(reflect.ValueOf(0)))
}

func TestRegistry_Definition(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	definition := NewTypedDefinition(registryTest{})
	registry.RegisterType(definition)

	require.Equal(definition, registry.Definition(reflect.ValueOf(registryTest{})))
	require.Equal(definition, registry.Definition(reflect.ValueOf(&registryTest{})))
	require.Nil(registry.Definition(reflect.ValueOf(registryNotFoundTest{})))
	require.Nil(registry.Definition(reflect.ValueOf(&registryNotFoundTest{})))
	require.Nil(registry.Definition(reflect.ValueOf(Validator(nil))))
	require.Nil(registry.Definition(reflect.ValueOf(nil)))
	require.Nil(registry.Definition(reflect.ValueOf(0)))
}
