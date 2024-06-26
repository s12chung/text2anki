// Package httptyped are helpers that check if the response is registered for it's type
package httptyped

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
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
			panic(fmt.Sprintf("httptyped.RegisterType() with non-Struct kind %v", typ.String()))
		}
		r.registeredTypes[typ] = true
	}
}

// HasType returns true if the type exists in the registry, also gives the name of the type
func (r *Registry) HasType(value any) bool {
	if value == nil {
		return false
	}
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

	for i := range typ.NumField() {
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

var errModelNil = errors.New("httptyped model is nil")

// PrepareModel checks if the type exists and prepares the model for serializing
func PrepareModel(model any) error {
	if model == nil {
		return errModelNil
	}
	if !HasType(model) {
		return fmt.Errorf("%v is not registered to httptyped", indirectTypeElement(reflect.TypeOf(model)).String())
	}

	if preparable, ok := model.(Preparable); ok {
		preparable.PrepareSerialize()
		return nil
	}
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := range val.Len() {
			elem := val.Index(i).Interface()
			preparable, ok := elem.(Preparable)
			if !ok {
				return nil
			}
			preparable.PrepareSerialize()
		}
	}
	return nil
}
