package types

import "reflect"

// TypeFunc func stuff
type TypeFunc func() interface{}

type Entry struct {
	Name         string
	Fullname     string
	Type         reflect.Type
	Factory      TypeFunc
	InternalType bool
}
