package httptyped

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s12chung/text2anki/pkg/util/httputil"
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
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fixture.CompareReadOrUpdate(t, path.Join(testName, tc.name+".json"), fixture.JSON(t, StructureMap(reflect.TypeOf(tc.str))))
		})
	}
}

type testObj struct {
	Val string `json:"val,omitempty"`
}

type invalidTestObj struct {
	Val string `json:"val,omitempty"`
}

func TestRespondTypedJSONWrap(t *testing.T) {
	DefaultRegistry.RegisterType(testObj{})

	var testVal string
	handlerFunc := RespondTypedJSONWrap(func(r *http.Request) (any, *httputil.HTTPError) {
		if r.Method == http.MethodPost {
			return nil, httputil.Error(http.StatusUnprocessableEntity, fmt.Errorf("not a GET"))
		}
		if r.Method == http.MethodPatch {
			return invalidTestObj{Val: testVal}, nil
		}
		if r.Method == http.MethodPut {
			return []invalidTestObj{{Val: testVal}}, nil
		}
		return testObj{Val: testVal}, nil
	})

	testCases := []struct {
		name         string
		method       string
		val          string
		status       int
		expectedBody string
	}{
		{name: "normal", val: "123", status: http.StatusOK, expectedBody: "{\"val\":\"123\"}\n"},
		{name: "err", method: http.MethodPost, status: http.StatusUnprocessableEntity,
			expectedBody: "{\"error\":\"not a GET\",\"code\":422,\"status_text\":\"Unprocessable Entity\"}\n"},
		{name: "not_registered", method: http.MethodPatch, status: http.StatusInternalServerError,
			expectedBody: "{\"error\":\"httptyped.invalidTestObj is not registered to httptyped\",\"code\":500,\"status_text\":\"Internal Server Error\"}\n"},
		{name: "slice_not_registered", method: http.MethodPut, status: http.StatusInternalServerError,
			expectedBody: "{\"error\":\"httptyped.invalidTestObj is not registered to httptyped\",\"code\":500,\"status_text\":\"Internal Server Error\"}\n"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			testVal = tc.val

			resp := httptest.NewRecorder()
			handlerFunc(resp, httptest.NewRequest(tc.method, "/", nil))

			require.Equal(tc.status, resp.Code)
			require.Equal(tc.expectedBody, resp.Body.String())
		})
	}
}
