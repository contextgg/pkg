package types

import (
	"reflect"
	"strings"
)

func TypeFactory(t reflect.Type) TypeFunc {
	return func() interface{} {
		return reflect.New(t).Interface()
	}
}

func GetElemType(source interface{}) reflect.Type {
	rawType := reflect.TypeOf(source)
	// source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}
	return rawType
}

func GetShortName(rawType reflect.Type) string {
	parts := strings.Split(rawType.String(), ".")
	return parts[1]
}

func GetTypeName(source interface{}) string {
	return GetShortName(GetElemType(source))
}
