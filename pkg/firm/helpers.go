package firm

import "reflect"

const nilName = "nil"

func typeName(value reflect.Value) string {
	if !value.IsValid() {
		return nilName
	}
	return indirect(value).Type().String()
}

func typeNameKey(value reflect.Value) ErrorKey {
	return ErrorKey(typeName(value))
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

func joinKeys(keys ...ErrorKey) ErrorKey {
	var key ErrorKey
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
