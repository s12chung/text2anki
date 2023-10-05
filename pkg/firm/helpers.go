package firm

import "reflect"

const nilName = "nil"

// TypeName returns the type name of the type
func TypeName(value reflect.Value) string {
	if !value.IsValid() {
		return nilName
	}
	return indirect(value).Type().String()
}

func indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	return value
}

func indirectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

const keySeparator = "."

func joinKeys[T ~string](keys ...T) T {
	var key T
	for _, v := range keys {
		if v == "" {
			continue
		}
		if key != "" {
			key += keySeparator
		}
		key += v
	}
	return key
}
