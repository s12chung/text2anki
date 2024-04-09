package httptyped

import (
	"encoding/json"
	"errors"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/test/fixture"
)

type Parent struct {
	Primitive int      `json:"primitive,omitempty"`
	Basic     Child    `json:"basic"`
	Pt        *Child   `json:"pt,omitempty"`
	Any       any      `json:"any,omitempty"`
	Array     []Child  `json:"array,omitempty"`
	ArrayPt   []*Child `json:"array_pt,omitempty"`
	NoTag     int
	SkipJSON  int `json:"-"`
	private   int //nolint:unused // For testing field
}

type Child struct {
	Str string `json:"str,omitempty"`
}

type Recurse struct {
	Recurse *Recurse `json:"recurse"`
	Layer   Layer    `json:"layer"`
}

type Layer struct {
	Recurse  *Recurse  `json:"recurse"`
	Recurses []Recurse `json:"recurses"`
}

type WithSerializedParent struct {
	Basic WithSerializedChild  `json:"basic"`
	Pt    *WithSerializedChild `json:"pt,omitempty"`
}

func (w WithSerializedParent) PrepareSerialize() { w.Pt.toSerialize = true }

type WithSerializedChild struct {
	toSerialize   bool
	NonSerialized string `json:"non_serialized,omitempty"`
}

type withSerializedChildAlias WithSerializedChild

type SerializedChild struct {
	Serialized string `json:"serialized,omitempty"`
}

func (w WithSerializedChild) SerializedEmpty() any { return SerializedChild{} }
func (w WithSerializedChild) ToSerialized() (any, error) {
	return SerializedChild{Serialized: w.NonSerialized + "_CEREAL"}, nil
}
func (w WithSerializedChild) MarshalJSON() ([]byte, error) {
	if !w.toSerialize {
		return json.Marshal(withSerializedChildAlias(w))
	}
	serialized, err := w.ToSerialized()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&serialized)
}

type WithSerializedParentPt struct {
	Basic WithSerializedChildPt  `json:"basic"`
	Pt    *WithSerializedChildPt `json:"pt,omitempty"`
}

func (w WithSerializedParentPt) PrepareSerialize() { w.Pt.toSerialize = true }

type WithSerializedChildPt struct {
	toSerialize   bool
	NonSerialized string `json:"non_serialized,omitempty"`
}

type SerializedChildPt struct {
	Serialized string `json:"serialized,omitempty"`
}

func (w *WithSerializedChildPt) SerializedEmpty() any { return SerializedChildPt{} }

func TestRegistry_RegisterType(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}

	require.False(registry.registeredTypes[reflect.TypeOf(Parent{})])
	require.False(registry.registeredTypes[reflect.TypeOf(Recurse{})])
	registry.RegisterType(Parent{}, Recurse{})
	require.True(registry.registeredTypes[reflect.TypeOf(Parent{})])
	require.True(registry.registeredTypes[reflect.TypeOf(Recurse{})])

	require.False(registry.registeredTypes[reflect.TypeOf(Child{})])
	registry.RegisterType(Child{})
	require.True(registry.registeredTypes[reflect.TypeOf(Child{})])

	require.Panics(func() {
		registry.RegisterType(1)
	})
}

func TestRegistry_HasType(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}

	registry.RegisterType(Parent{})
	require.True(registry.HasType(Parent{}))
	require.False(registry.HasType(Child{}))
}

func TestRegistry_Types(t *testing.T) {
	require := require.New(t)

	registry := &Registry{}
	registry.RegisterType(Parent{})
	registry.RegisterType(Recurse{})
	require.ElementsMatch([]reflect.Type{reflect.TypeOf(Parent{}), reflect.TypeOf(Recurse{})}, registry.Types())
}

func TestStructureMap(t *testing.T) {
	testName := "TestStructureMap"

	testCases := []struct {
		name string
		str  any
	}{
		{name: "Parent", str: Parent{}},
		{name: "Child", str: Child{}},
		{name: "Recurse", str: Recurse{}},
		{name: "WithSerialized", str: WithSerializedParent{}},
		{name: "WithSerializedPt", str: WithSerializedParentPt{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, tc.name), StructureMap(reflect.TypeOf(tc.str)))
		})
	}
}

type testObj struct {
	Val string `json:"val,omitempty"`
}

type invalidTestObj struct {
	Val string `json:"val,omitempty"`
}

func TestPrepareModel(t *testing.T) {
	testName := "TestPrepareModel"
	DefaultRegistry.RegisterType(testObj{})
	DefaultRegistry.RegisterType(WithSerializedParent{})
	notRegisteredErr := errors.New("httptyped.invalidTestObj is not registered to httptyped")

	testCases := []struct {
		name  string
		model any
		err   error
	}{
		{name: "toSerialize", model: newPrepareModel()},
		{name: "toSerialize_slice", model: []WithSerializedParent{newPrepareModel(), newPrepareModel()}},
		{name: "not_registered", model: invalidTestObj{Val: "123"}, err: notRegisteredErr},
		{name: "slice_not_registered", model: []invalidTestObj{{Val: "123"}}, err: notRegisteredErr},
		{name: "nil", err: errModelNil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			err := PrepareModel(tc.model)
			if tc.err != nil {
				require.Equal(tc.err, err)
				return
			}
			require.NoError(err)
			fixture.CompareReadOrUpdateJSON(t, path.Join(testName, tc.name), tc.model)
		})
	}
}

func newPrepareModel() WithSerializedParent {
	return WithSerializedParent{
		Basic: WithSerializedChild{NonSerialized: "Basic"},
		Pt:    &WithSerializedChild{NonSerialized: "Pt"},
	}
}
