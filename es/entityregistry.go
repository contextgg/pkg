package es

import (
	"errors"
	"fmt"
	"sync"

	"github.com/contextgg/pkg/types"
)

var ErrEntityNotFound = fmt.Errorf("Entity not found")

// EntityRegistry stores the handlers for commands
type EntityRegistry interface {
	GetOptions(entityName string) (EntityOptions, error)
	SetEntity(entityType interface{}, opts ...EntityOption) error
}

type entityRegistry struct {
	sync.RWMutex
	registry map[string]EntityOptions
}

func (r *entityRegistry) SetEntity(entityType interface{}, opts ...EntityOption) error {
	r.Lock()
	defer r.Unlock()

	options := NewEntityOptions(opts)
	if options.Factory == nil {
		return errors.New("You need to supply a factory method")
	}

	name := types.GetTypeName(entityType)
	r.registry[name] = options
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
	}
}
