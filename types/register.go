package types

import (
	"fmt"
	"reflect"
	"strings"
)

type Register interface {
	GetAll() map[string]*TypeData
	GetTypeData(name string) (*TypeData, bool)
	SetTypeDataByFactory(fn TypeFunc, internalType bool) *TypeData
	SetTypeData(t interface{}, internalType bool) *TypeData
	UnmarshalByName(name string, raw []byte, legacy bool) (interface{}, error)
}

type register struct {
	names map[string]*TypeData
}

// SetTypeDataByFactory set the type
func (r *register) SetTypeDataByFactory(fn TypeFunc, internalType bool) *TypeData {
	rawType, name := GetTypeName(fn())
	return r.set(name, rawType, fn, internalType)
}

// SetTypeData set the type
func (r *register) SetTypeData(t interface{}, internalType bool) *TypeData {
	rawType, name := GetTypeName(t)
	fn := func() interface{} {
		return reflect.New(rawType).Interface()
	}
	return r.set(name, rawType, fn, internalType)
}

// GetTypeData will try resolve type data by a name
func (r *register) GetTypeData(name string) (*TypeData, bool) {
	lower := strings.ToLower(name)
	if data, ok := r.names[lower]; ok {
		return data, true
	}
	return nil, false
}

func (r *register) GetAll() map[string]*TypeData {
	return r.names
}

func (r *register) UnmarshalByName(name string, raw []byte, legacy bool) (interface{}, error) {
	td, ok := r.GetTypeData(name)
	if !ok {
		return nil, fmt.Errorf("No type found with name %s", name)
	}
	data := td.Factory()

	return Unmarshal(data, raw, legacy)
}

func (r *register) set(name string, rawType reflect.Type, fn TypeFunc, internalType bool) *TypeData {
	lower := strings.ToLower(name)
	data := &TypeData{
		Name:         name,
		Factory:      fn,
		Type:         rawType,
		InternalType: internalType,
	}
	r.names[lower] = data
	return data
}

func NewRegistry() Register {
	return &register{
		names: make(map[string]*TypeData),
	}
}
