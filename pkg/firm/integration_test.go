package firm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type parent struct {
	Primitive               int
	Basic                   child
	Pt                      *child
	Any                     any
	PrimitiveEmptyValidates int
	BasicEmptyValidates     child
	PtEmptyValidates        *child
	AnyEmptyValidates       any
	PrimitiveNoValidates    int
	BasicNoValidates        child
	PtNoValidates           *child
	AnyNoValidates          any
}

type child struct {
	Validates   string
	NoValidates string
}

func fullParent() parent {
	fullChild := child{Validates: "child validates", NoValidates: "no validates"}

	fc1 := fullChild
	fc2 := fullChild
	fc3 := fullChild
	return parent{
		// validate field + child
		Primitive: 1, Basic: fullChild, Pt: &fc1, Any: fullChild,
		// validate child
		PrimitiveEmptyValidates: 1, BasicEmptyValidates: fullChild, PtEmptyValidates: &fc2, AnyEmptyValidates: fullChild,
		// validate none
		PrimitiveNoValidates: 1, BasicNoValidates: fullChild, PtNoValidates: &fc3, AnyNoValidates: fullChild,
	}
}

type topLevelValidates struct {
	Primitive  int
	Primitive2 int
}

type unregistered struct{}

var testRegistry = &Registry{}

func init() {
	testRegistry.RegisterType(
		NewTypedDefinition(parent{}).
			Validates(RuleMap{
				"Primitive":               {testPresence{}},
				"Basic":                   {testPresence{}},
				"Pt":                      {testPresence{}},
				"Any":                     {testPresence{}},
				"PrimitiveEmptyValidates": {},
				"BasicEmptyValidates":     {},
				"PtEmptyValidates":        {},
				"AnyEmptyValidates":       {},
			}))
	testRegistry.RegisterType(
		NewTypedDefinition(child{}).
			Validates(RuleMap{
				"Validates": {testPresence{}},
			}))
	testRegistry.RegisterType(
		NewTypedDefinition(topLevelValidates{}).
			ValidatesTopLevel(testPresence{}))
}

type integrationTestCase struct {
	name    string
	isValid bool
	f       func() parent
	anyF    func() any
}

var integrationAnyTestCases = []integrationTestCase{
	//
	// Any
	//
	{name: "Data___int_raw", isValid: false, anyF: func() any {
		return 1
	}},
	{name: "Data___int_pt", isValid: false, anyF: func() any {
		i := 1
		return &i
	}},
	{name: "Data___unregistered_raw", isValid: false, anyF: func() any {
		return unregistered{}
	}},
	{name: "Data___unregistered_pt", isValid: false, anyF: func() any {
		return &unregistered{}
	}},
	{name: "Data___nil_raw", isValid: false, anyF: func() any {
		return nil
	}},
	{name: "Data___nil_pt", isValid: false, anyF: func() any {
		var i any
		return &i
	}},
	{name: "Data___topLevelValidates_full", isValid: true, anyF: func() any {
		return topLevelValidates{Primitive: 1, Primitive2: 2}
	}},
	{name: "Data___topLevelValidates_half_raw", isValid: true, anyF: func() any {
		return topLevelValidates{Primitive: 1}
	}},
	{name: "Data___topLevelValidates_half_pt", isValid: true, anyF: func() any {
		return &topLevelValidates{Primitive: 1}
	}},
	{name: "Data___topLevelValidates_empty_raw", isValid: false, anyF: func() any {
		return topLevelValidates{}
	}},
	{name: "Data___topLevelValidates_empty_pt", isValid: false, anyF: func() any {
		return &topLevelValidates{}
	}},
	{name: "Empty", isValid: false, f: func() parent {
		return parent{}
	}},
	{name: "Empty___any_raw", isValid: false, anyF: func() any {
		return parent{}
	}},
	{name: "Empty___any_pt", isValid: false, anyF: func() any {
		return &parent{}
	}},
	{name: "Full___any_raw", isValid: true, anyF: func() any {
		return fullParent()
	}},
	{name: "Full___any_pt", isValid: true, anyF: func() any {
		full := fullParent()
		return &full
	}},
}

func TestIntegration(t *testing.T) {
	integrationTestCases := make([]integrationTestCase, len(structValidatorTestCases))
	for i, v := range structValidatorTestCases {
		integrationTestCases[i] = integrationTestCase{
			name:    v.name,
			isValid: len(v.errorKeys) == 0,
			f:       v.f,
		}
	}

	for _, tc := range append(integrationAnyTestCases, integrationTestCases...) {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			if tc.f != nil {
				data := tc.f()
				require.Equal(tc.isValid, testRegistry.Validate(data).IsValid())
				require.Equal(tc.isValid, testRegistry.Validate(&data).IsValid())
				return
			}
			require.Equal(tc.isValid, testRegistry.Validate(tc.anyF()).IsValid())
		})
	}
}
