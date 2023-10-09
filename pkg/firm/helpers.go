package firm

import "reflect"

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
