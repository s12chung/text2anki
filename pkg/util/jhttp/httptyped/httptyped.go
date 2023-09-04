// Package httptyped are helpers that check if the response is registered for it's type
package httptyped

import (
	"fmt"
	"net/http"
	"reflect"
	"unicode"

	"github.com/s12chung/text2anki/pkg/util/jhttp"
)

// RegisterType registers the type to the DefaultRegistry
var RegisterType = DefaultRegistry.RegisterType

// HasType checks the type is registered in the DefaultRegistry
var HasType = DefaultRegistry.HasType

// Types returns the registered types in the DefaultRegistry
var Types = DefaultRegistry.Types

// DefaultRegistry is the registry used for global functions
var DefaultRegistry = &Registry{}

// Registry contains the types registered and their formats
type Registry struct {
	registeredTypes map[reflect.Type]bool
}

// RegisterType adds the type to the registry
func (r *Registry) RegisterType(values ...any) {
	if r.registeredTypes == nil {
		r.registeredTypes = map[reflect.Type]bool{}
	}

	for _, value := range values {
		typ := indirectType(reflect.TypeOf(value))
		if typ.Kind() != reflect.Struct {
			panic(fmt.Sprintf("RegisterType() with non-Struct kind %v", typ.String()))
		}
		r.registeredTypes[typ] = true
	}
}

// HasType returns true if the type exists in the registry, also gives the name of the type
func (r *Registry) HasType(value any) bool {
	return r.registeredTypes[indirectTypeElement(reflect.TypeOf(value))]
}

// Types returns the types in the registry
func (r *Registry) Types() []reflect.Type {
	types := make([]reflect.Type, len(r.registeredTypes))
	i := 0
	for k := range r.registeredTypes {
		types[i] = k
		i++
	}
	return types
}

func indirectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func indirectTypeElement(typ reflect.Type) reflect.Type {
	typ = indirectType(typ)
	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		typ = typ.Elem()
	}
	typ = indirectType(typ)
	return typ
}

const serializedEmptyFunctionName = "SerializedEmpty"

// HasSerialized are models that has a Serialized version of the model
type HasSerialized interface {
	SerializedEmpty() any
}

var hasSerializedType = reflect.TypeOf((*HasSerialized)(nil)).Elem()

func serializedType(typ reflect.Type) reflect.Type {
	isPointer := typ.Kind() == reflect.Pointer
	typ = indirectSerializedType(indirectType(typ))
	if isPointer {
		typ = reflect.PointerTo(typ)
	}
	return typ
}

func indirectSerializedType(typ reflect.Type) reflect.Type {
	var method reflect.Value
	if typ.Implements(hasSerializedType) {
		method = reflect.New(typ).Elem().MethodByName(serializedEmptyFunctionName)
	} else {
		typPointer := reflect.PointerTo(typ)
		if typPointer.Implements(hasSerializedType) {
			method = reflect.New(typPointer).Elem().MethodByName(serializedEmptyFunctionName)
		}
	}
	if !method.IsValid() {
		return typ
	}
	return method.Call(nil)[0].Elem().Type()
}

// StructureMap returns a map of the structure of the type
func StructureMap(typ reflect.Type) map[string]map[string]string {
	m := map[string]map[string]string{}
	structureMap(typ, m, map[reflect.Type]bool{})
	return m
}

func structureMap(typ reflect.Type, m map[string]map[string]string, handledTypeMap map[reflect.Type]bool) {
	currentTypeMap := map[string]string{}
	typ = serializedType(typ)

	m[typ.String()] = currentTypeMap
	handledTypeMap[typ] = true

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if unicode.IsLower([]rune(field.Name)[0]) {
			continue
		}
		jsonKey := field.Tag.Get("json")
		if jsonKey == "" {
			jsonKey = field.Name
		}
		if jsonKey == "-" {
			continue
		}
		currentTypeMap[jsonKey] = serializedType(field.Type).String()

		fieldType := indirectTypeElement(field.Type)
		if !handledTypeMap[fieldType] && fieldType.Kind() == reflect.Struct {
			structureMap(fieldType, m, handledTypeMap)
		}
	}
}

// Preparable is a model that has a Serialized version of itself or embedded
type Preparable interface {
	PrepareSerialize()
}

// TypedWrap wraps around jhttp.ResponseJSONWrap, but also checks the type of the response beforehand
func TypedWrap(f jhttp.ResponseJSONWrapFunc) jhttp.ResponseJSONWrapFunc {
	return func(r *http.Request) (any, *jhttp.HTTPError) {
		resp, httpError := f(r)
		if httpError != nil {
			return resp, httpError
		}
		if !HasType(resp) {
			return nil, jhttp.Error(http.StatusInternalServerError,
				fmt.Errorf("%v is not registered to httptyped", indirectTypeElement(reflect.TypeOf(resp))))
		}
		if preparable, ok := resp.(Preparable); ok {
			preparable.PrepareSerialize()
		}
		return resp, nil
	}
}
