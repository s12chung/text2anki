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

//nolint:funlen,maintidx // it's a big testing function
func TestIntegration(t *testing.T) {
	type testCase struct {
		name     string
		expected bool
		f        func() parent
		anyF     func() any
	}
	tcs := []testCase{
		//
		// Any
		//
		{name: "Data___int_raw", expected: false, anyF: func() any {
			return 1
		}},
		{name: "Data___int_pt", expected: false, anyF: func() any {
			i := 1
			return &i
		}},
		{name: "Data___unregistered_raw", expected: false, anyF: func() any {
			return unregistered{}
		}},
		{name: "Data___unregistered_pt", expected: false, anyF: func() any {
			return &unregistered{}
		}},
		{name: "Data___nil_raw", expected: false, anyF: func() any {
			return nil
		}},
		{name: "Data___nil_pt", expected: false, anyF: func() any {
			var i any
			return &i
		}},
		{name: "Data___topLevelValidates_full", expected: true, anyF: func() any {
			return topLevelValidates{Primitive: 1, Primitive2: 2}
		}},
		{name: "Data___topLevelValidates_half_raw", expected: true, anyF: func() any {
			return topLevelValidates{Primitive: 1}
		}},
		{name: "Data___topLevelValidates_half_pt", expected: true, anyF: func() any {
			return &topLevelValidates{Primitive: 1}
		}},
		{name: "Data___topLevelValidates_empty_raw", expected: false, anyF: func() any {
			return topLevelValidates{}
		}},
		{name: "Data___topLevelValidates_empty_pt", expected: false, anyF: func() any {
			return &topLevelValidates{}
		}},
		{name: "Empty", expected: false, f: func() parent {
			return parent{}
		}},
		{name: "Empty___any_raw", expected: false, anyF: func() any {
			return parent{}
		}},
		{name: "Empty___any_pt", expected: false, anyF: func() any {
			return &parent{}
		}},
		{name: "Full___any_raw", expected: true, anyF: func() any {
			return fullParent()
		}},
		{name: "Full___any_pt", expected: true, anyF: func() any {
			full := fullParent()
			return &full
		}},

		//
		// Full
		//
		{name: "Full", expected: true, f: fullParent},

		//
		// Primitive
		//
		{name: "Primitive___zero", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Primitive = 0
			return changeParent
		}},

		//
		// Basic
		//
		{name: "Basic___child_validates_ok", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.Basic.NoValidates = ""
			return changeParent
		}},
		{name: "Basic___child_validates_zero", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Basic.Validates = ""
			return changeParent
		}},
		{name: "Basic___child_empty", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Basic = child{}
			return changeParent
		}},

		//
		// Pt
		//
		{name: "Pt___child_validates_ok", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.Pt.NoValidates = ""
			return changeParent
		}},
		{name: "Pt___child_validates_zero", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Pt.Validates = ""
			return changeParent
		}},
		{name: "Pt___child_empty", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Pt = &child{}
			return changeParent
		}},
		{name: "Pt___nil", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Pt = nil
			return changeParent
		}},

		//
		// Any
		//
		{name: "Any___child_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.Any = child{}
			return changeParent
		}},
		{name: "Any___child_pointer_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.Any = &child{}
			return changeParent
		}},
		{name: "Any___nil", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.Any = nil
			return changeParent
		}},

		//
		// PrimitiveEmptyValidates
		//
		{name: "PrimitiveEmptyValidates___zero", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PrimitiveEmptyValidates = 0
			return changeParent
		}},

		//
		// BasicEmptyValidates
		//
		{name: "BasicEmptyValidates___child_validates_ok", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.BasicEmptyValidates.NoValidates = ""
			return changeParent
		}},
		{name: "BasicEmptyValidates___child_validates_zero", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.BasicEmptyValidates.Validates = ""
			return changeParent
		}},
		{name: "BasicEmptyValidates____child_empty", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.BasicEmptyValidates = child{}
			return changeParent
		}},

		//
		// PtEmptyValidates
		//
		{name: "PtEmptyValidates___child_validates_ok", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PtEmptyValidates.NoValidates = ""
			return changeParent
		}},
		{name: "PtEmptyValidates___child_validates_zero", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.PtEmptyValidates.Validates = ""
			return changeParent
		}},
		{name: "PtEmptyValidates___child_empty", expected: false, f: func() parent {
			changeParent := fullParent()
			changeParent.PtEmptyValidates = &child{}
			return changeParent
		}},
		{name: "PtEmptyValidates___nil", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PtEmptyValidates = nil
			return changeParent
		}},

		//
		// AnyEmptyValidates
		//
		{name: "AnyEmptyValidates___child_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.AnyEmptyValidates = child{}
			return changeParent
		}},
		{name: "AnyEmptyValidates___child_pointer_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.AnyEmptyValidates = &child{}
			return changeParent
		}},
		{name: "AnyEmptyValidates___nil", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.AnyEmptyValidates = nil
			return changeParent
		}},

		//
		// PrimitiveNoValidates
		//
		{name: "PrimitiveNoValidates___zero", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PrimitiveNoValidates = 0
			return changeParent
		}},

		//
		// BasicNoValidates
		//
		{name: "BasicNoValidates___child_validates_ok", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.BasicNoValidates.NoValidates = ""
			return changeParent
		}},
		{name: "BasicNoValidates___child_validates_zero", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.BasicNoValidates.Validates = ""
			return changeParent
		}},
		{name: "BasicNoValidates____child_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.BasicNoValidates = child{}
			return changeParent
		}},

		//
		// PtNoValidates
		//
		{name: "PtNoValidates___child_validates_ok", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PtNoValidates.NoValidates = ""
			return changeParent
		}},
		{name: "PtNoValidates___child_validates_zero", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PtNoValidates.Validates = ""
			return changeParent
		}},
		{name: "PtNoValidates___child_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PtNoValidates = &child{}
			return changeParent
		}},
		{name: "PtNoValidates___nil", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.PtNoValidates = nil
			return changeParent
		}},

		//
		// AnyNoValidates
		//
		{name: "AnyNoValidates___child_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.AnyNoValidates = child{}
			return changeParent
		}},
		{name: "AnyNoValidates___child_pointer_empty", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.AnyNoValidates = &child{}
			return changeParent
		}},
		{name: "AnyNoValidates___nil", expected: true, f: func() parent {
			changeParent := fullParent()
			changeParent.AnyNoValidates = nil
			return changeParent
		}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			if tc.f != nil {
				data := tc.f()
				require.Equal(tc.expected, testRegistry.Validate(data).IsValid())
				require.Equal(tc.expected, testRegistry.Validate(&data).IsValid())
				return
			}
			require.Equal(tc.expected, testRegistry.Validate(tc.anyF()).IsValid())
		})
	}
}
