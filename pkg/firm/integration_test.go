package firm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type validTestCase struct {
	Primitive         int             `validates:"presence"`
	Basic             validTestChild  `validates:"presence"`
	Pt                *validTestChild `validates:"presence"`
	Any               any             `validates:"presence"`
	PrimitiveEmptyTag int             `validates:""`
	BasicEmptyTag     validTestChild  `validates:""`
	PtEmptyTag        *validTestChild `validates:""`
	AnyEmptyTag       any             `validates:""`
	PrimitiveNoTag    int
	BasicNoTag        validTestChild
	PtNoTag           *validTestChild
	AnyNoTag          any
}

type validTestChild struct {
	Tag   string `validates:"presence"`
	NoTag string
}

//nolint:funlen,maintidx // it's a big testing function
func TestValidator_IsValid(t *testing.T) {
	fullChild := validTestChild{Tag: "fullCase", NoTag: "no tag"}
	fullCase := func() validTestCase {
		pt := fullChild
		ptEmptyTag := fullChild
		ptNoTag := fullChild
		return validTestCase{
			// validate validTestCase + validTestChild
			Primitive: 1, Basic: fullChild, Pt: &pt, Any: fullChild,
			// validate validTestChild
			PrimitiveEmptyTag: 1, BasicEmptyTag: fullChild, PtEmptyTag: &ptEmptyTag, AnyEmptyTag: fullChild,
			// validate none
			PrimitiveNoTag: 1, BasicNoTag: fullChild, PtNoTag: &ptNoTag, AnyNoTag: fullChild,
		}
	}

	type testCase struct {
		name         string
		expected     bool
		testCaseFunc func() validTestCase
		anyFunc      func() any
	}
	tcs := []testCase{
		{name: "Data___int", expected: false, anyFunc: func() any {
			return 1
		}},
		{name: "Full", expected: true, testCaseFunc: func() validTestCase {
			return fullCase()
		}},

		//
		// Primitive
		//
		{name: "Primitive___zero", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Primitive = 0
			return childCase
		}},

		//
		// Basic
		//
		{name: "Basic___child_tagged", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Basic.NoTag = ""
			return childCase
		}},
		{name: "Basic___child_tag_zero", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Basic.Tag = ""
			return childCase
		}},
		{name: "Basic___child_empty", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Basic = validTestChild{}
			return childCase
		}},

		//
		// Pt
		//
		{name: "Pt___child_tagged", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Pt.NoTag = ""
			return childCase
		}},
		{name: "Pt___child_tag_zero", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Pt.Tag = ""
			return childCase
		}},
		{name: "Pt___child_empty", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Pt = &validTestChild{}
			return childCase
		}},
		{name: "Pt___nil", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Pt = nil
			return childCase
		}},

		//
		// Any
		//
		{name: "Any___child_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Any = validTestChild{}
			return childCase
		}},
		{name: "Any___child_pointer_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Any = &validTestChild{}
			return childCase
		}},
		{name: "Any___nil", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.Any = nil
			return childCase
		}},

		//
		// PrimitiveEmptyTag
		//
		{name: "PrimitiveEmptyTag___zero", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PrimitiveEmptyTag = 0
			return childCase
		}},

		//
		// BasicEmptyTag
		//
		{name: "BasicEmptyTag___child_tagged", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.BasicEmptyTag.NoTag = ""
			return childCase
		}},
		{name: "BasicEmptyTag___child_tag_zero", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.BasicEmptyTag.Tag = ""
			return childCase
		}},
		{name: "BasicEmptyTag____child_empty", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.BasicEmptyTag = validTestChild{}
			return childCase
		}},

		//
		// PtEmptyTag
		//
		{name: "PtEmptyTag___child_tagged", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtEmptyTag.NoTag = ""
			return childCase
		}},
		{name: "PtEmptyTag___child_tag_zero", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtEmptyTag.Tag = ""
			return childCase
		}},
		{name: "PtEmptyTag___child_empty", expected: false, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtEmptyTag = &validTestChild{}
			return childCase
		}},
		{name: "PtEmptyTag___nil", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtEmptyTag = nil
			return childCase
		}},

		//
		// AnyEmptyTag
		//
		{name: "AnyEmptyTag___child_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.AnyEmptyTag = validTestChild{}
			return childCase
		}},
		{name: "AnyEmptyTag___child_pointer_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.AnyEmptyTag = &validTestChild{}
			return childCase
		}},
		{name: "AnyEmptyTag___nil", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.AnyEmptyTag = nil
			return childCase
		}},

		//
		// PrimitiveNoTag
		//
		{name: "PrimitiveNoTag___zero", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PrimitiveNoTag = 0
			return childCase
		}},

		//
		// BasicNoTag
		//
		{name: "BasicNoTag___child_tagged", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.BasicNoTag.NoTag = ""
			return childCase
		}},
		{name: "BasicNoTag___child_tag_zero", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.BasicNoTag.Tag = ""
			return childCase
		}},
		{name: "BasicNoTag____child_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.BasicNoTag = validTestChild{}
			return childCase
		}},

		//
		// PtNoTag
		//
		{name: "PtNoTag___child_tagged", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtNoTag.NoTag = ""
			return childCase
		}},
		{name: "PtNoTag___child_tag_zero", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtNoTag.Tag = ""
			return childCase
		}},
		{name: "PtNoTag___child_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtNoTag = &validTestChild{}
			return childCase
		}},
		{name: "PtNoTag___nil", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.PtNoTag = nil
			return childCase
		}},

		//
		// AnyNoTag
		//
		{name: "AnyNoTag___child_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.AnyNoTag = validTestChild{}
			return childCase
		}},
		{name: "AnyNoTag___child_pointer_empty", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.AnyNoTag = &validTestChild{}
			return childCase
		}},
		{name: "AnyNoTag___nil", expected: true, testCaseFunc: func() validTestCase {
			childCase := fullCase()
			childCase.AnyNoTag = nil
			return childCase
		}},
	}

	runFunc := func(t *testing.T, tc testCase) {
		require := require.New(t)
		if tc.testCaseFunc != nil {
			require.Equal(tc.expected, New(tc.testCaseFunc()).IsValid())
			return
		}
		require.Equal(tc.expected, New(tc.anyFunc()).IsValid())
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			runFunc(t, tc)
		})
	}
	for _, tc := range tcs {
		t.Run("Pointer_"+tc.name, func(t *testing.T) {
			runFunc(t, tc)
		})
	}
}
