package firm

import (
	"reflect"
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
	registry.RegisterType(NewDefinition(registryTestParent{}).
		Validates(RuleMap{"Child": {}}).
		ValidatesTopLevel(testPresence{}))

	typeToValidator := map[reflect.Type]*ValueValidator{
		reflect.TypeOf(registryTestParent{}): {Rules: []Rule{testPresence{}, &StructValidator{RuleMap: map[string]*[]Rule{
			"Child": {},
		}}}},
	}
	require.Equal(typeToValidator, registry.typeToValidator)
	require.Equal(map[reflect.Type][]*[]Rule{
		reflect.TypeOf(registryTestChild{}): {{}},
	}, registry.unregisteredTypeReferences)

	registry.RegisterType(NewDefinition(registryTestChild{}))

	typeToValidator[reflect.TypeOf(registryTestParent{})] = &ValueValidator{Rules: []Rule{testPresence{},
		&StructValidator{RuleMap: map[string]*[]Rule{
			"Child": {&ValueValidator{Rules: []Rule{&StructValidator{RuleMap: map[string]*[]Rule{}}}}},
		}}}}
	typeToValidator[reflect.TypeOf(registryTestChild{})] = &ValueValidator{Rules: []Rule{&StructValidator{RuleMap: map[string]*[]Rule{}}}}
	require.Equal(typeToValidator, registry.typeToValidator)
	require.Equal(map[reflect.Type][]*[]Rule{}, registry.unregisteredTypeReferences)

	require.Panics(func() {
		registry.RegisterType(NewDefinition(registryTestParent{}).ValidatesTopLevel(testPresence{}))
	})
}

// nolint:funlen // a bunch of test cases
func TestRegistry_Validate(t *testing.T) {
	type testCase struct {
		name       string
		definition *Definition
		data       any

		expectedKeySuffix ErrorKey
		err               *TemplatedError
	}
	tcs := []testCase{
		{
			name:              "top_level",
			definition:        NewDefinition(registryTestParent{}).ValidatesTopLevel(testPresence{}),
			data:              registryTestParent{},
			expectedKeySuffix: testPresenceKey,
			err:               errTest,
		},
		{
			name: "field_Primitive",
			definition: NewDefinition(registryTestParent{}).
				Validates(RuleMap{
					"Primitive": {testPresence{}},
				}),
			data:              registryTestParent{},
			expectedKeySuffix: "Primitive.testPresence",
			err:               errTest,
		},
		{
			name:              "not_found",
			definition:        NewDefinition(registryTestParent{}).ValidatesTopLevel(testPresence{}),
			data:              registryNotFoundTest{},
			expectedKeySuffix: notFoundRuleErrorKey,
			err:               notFoundRuleError(reflect.ValueOf(registryNotFoundTest{})),
		},
		{
			name: "not_found_field_Primitive",
			definition: NewDefinition(registryTestParent{}).
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
			registry := &Registry{}
			registry.RegisterType(tc.definition)

			testValidates(t, registry, tc.data, tc.err, tc.expectedKeySuffix)
		})
	}
}

func TestRegistry_DefaultedValidator(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	registry.RegisterType(NewDefinition(registryTestParent{}).ValidatesTopLevel(testPresence{}))

	structValidator := NewStructValidator(nil)
	expected := NewValueValidator(testPresence{}, &structValidator)
	require.Equal(&expected, registry.DefaultedValidator(reflect.ValueOf(registryTestParent{})))

	notFoundValue := reflect.ValueOf(nil)
	require.Equal(DefaultValidator, registry.DefaultedValidator(notFoundValue))

	registry.DefaultValidator = ValueValidator{}
	require.Equal(Validator(ValueValidator{}), registry.DefaultedValidator(notFoundValue))
}

func TestRegistry_Validator(t *testing.T) {
	testRegistryValidatorF(t, func(registry *Registry, data any) any {
		return registry.Validator(reflect.ValueOf(data))
	})
}

func TestRegistry_ValidatorForType(t *testing.T) {
	testRegistryValidatorF(t, func(registry *Registry, data any) any {
		return registry.ValidatorForType(reflect.TypeOf(data))
	})
}

func testRegistryValidatorF(t *testing.T, f func(registry *Registry, data any) any) {
	require := require.New(t)

	registry := &Registry{}
	registry.RegisterType(NewDefinition(registryTestParent{}).ValidatesTopLevel(testPresence{}))

	structValidator := NewStructValidator(nil)
	expected := NewValueValidator(testPresence{}, &structValidator)
	require.Equal(&expected, f(registry, registryTestParent{}))
	require.Equal(&expected, f(registry, &registryTestParent{}))
	require.Nil(f(registry, registryNotFoundTest{}))
	require.Nil(f(registry, &registryNotFoundTest{}))
	require.Nil(f(registry, Validator(nil)))
	require.Nil(f(registry, nil))
	require.Nil(f(registry, 0))
}
