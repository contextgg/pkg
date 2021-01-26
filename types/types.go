package types

import (
	"reflect"
)

// TypeData holds information regarding a type
type TypeData struct {
	Name         string
	Type         reflect.Type
	Factory      TypeFunc
	InternalType bool
}

// TypeFunc func stuff
type TypeFunc func() interface{}
