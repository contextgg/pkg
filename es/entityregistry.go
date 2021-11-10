package es

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/contextgg/pkg/types"
)

var ErrEntityNotFound = fmt.Errorf("Entity not found")

func entityOptions(options []EntityOption) EntityOptions {
	// set defaults.
	o := EntityOptions{
		Revision:       "rev1",
		Project:        true,
		MinVersionDiff: 0,
	}

	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	return o
}

// EntityRegistry stores the handlers for commands
type EntityRegistry interface {
	GetOptions(entityName string) (EntityOptions, error)
	SetEntity(entityType EntityType, opts ...EntityOption) error
}

type entityRegistry struct {
	sync.RWMutex
	registry map[string]EntityOptions
	types    map[string]reflect.Type
}

func (r *entityRegistry) SetEntity(entityType EntityType, opts ...EntityOption) error {
	r.Lock()
	defer r.Unlock()

	options := entityOptions(opts)
	if options.Factory == nil {
		return errors.New("You need to supply a factory method")
	}

	rawType, name := types.GetTypeName(entityType)
	r.registry[name] = options
	r.types[name] = rawType
	return nil
}

func (r *entityRegistry) GetOptions(entityName string) (EntityOptions, error) {
	opt, ok := r.registry[entityName]
	if !ok {
		return EntityOptions{}, ErrEntityNotFound
	}

	return opt, nil
}

// NewCommandRegistry creates a new CommandRegistry
func NewEntityRegistry() EntityRegistry {
	return &entityRegistry{
		registry: make(map[string]EntityOptions),
		types:    make(map[string]reflect.Type),
	}
}
