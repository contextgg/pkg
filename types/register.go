package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInvalidName       = fmt.Errorf("Invalid name")
	ErrInvalidFullName   = fmt.Errorf("Invalid fullname")
	ErrInvalidType       = fmt.Errorf("Invalid type")
	ErrAlreadyRegistered = fmt.Errorf("Already registered")
	ErrNotFound          = fmt.Errorf("Entry not found")
)

func condenseErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}
	err := errs[0]
	for _, e := range errs[1:] {
		err = errors.Wrap(err, e.Error())
	}
	return err
}

type Registry interface {
	All() []*Entry
	Add(options ...TypeOption) (*Entry, error)
	GetFirstByNames(names []string) (*Entry, bool)
	GetByName(name string) (*Entry, bool)
	GetByType(t reflect.Type) (*Entry, bool)
}

type registry struct {
	names map[string]*Entry
	types map[reflect.Type]*Entry
}

func (r *registry) All() []*Entry {
	var all []*Entry
	for _, value := range r.names {
		all = append(all, value)
	}
	return all
}

func (r *registry) Add(options ...TypeOption) (*Entry, error) {
	opts := &Entry{}
	for _, o := range options {
		o(opts)
	}

	// validate!
	var errors []error
	if len(opts.Name) == 0 {
		errors = append(errors, ErrInvalidName)
	}
	if len(opts.Fullname) == 0 {
		errors = append(errors, ErrInvalidFullName)
	}
	if opts.Type == nil {
		errors = append(errors, ErrInvalidType)
	}
	if len(errors) > 0 {
		return nil, condenseErrors(errors)
	}

	// do we need a factory?
	if opts.Factory == nil {
		opts.Factory = typeFactory(opts.Type)
	}

	// do we already have one registered?
	lower := strings.ToLower(opts.Name)
	if _, ok := r.names[lower]; ok {
		return nil, ErrAlreadyRegistered
	}

	r.names[lower] = opts
	r.types[opts.Type] = opts
	return opts, nil
}
func (r *registry) GetFirstByNames(names []string) (*Entry, bool) {
	for _, name := range names {
		if entry, ok := r.GetByName(name); ok {
			return entry, true
		}
	}
	return nil, false
}
func (r *registry) GetByName(name string) (*Entry, bool) {
	parts := strings.Split(name, ".")
	out := parts[len(parts)-1]
	lower := strings.ToLower(out)

	entry, ok := r.names[lower]
	return entry, ok
}
func (r *registry) GetByType(t reflect.Type) (*Entry, bool) {
	entry, ok := r.types[t]
	return entry, ok
}

func NewRegistry() Registry {
	return &registry{
		names: make(map[string]*Entry),
		types: make(map[reflect.Type]*Entry),
	}
}
