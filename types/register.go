package types

import (
	"fmt"
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
	Add(...*Entry) error
	GetByName(name string) (*Entry, bool)
}

type registry struct {
	names map[string]*Entry
}

func (r *registry) Add(entries ...*Entry) error {
	for _, e := range entries {
		// register all names.
		for _, name := range e.Names {
			// do we already have one registered?
			lower := strings.ToLower(name)
			if _, ok := r.names[lower]; ok {
				return ErrAlreadyRegistered
			}

			r.names[lower] = e
		}
	}
	return nil
}
func (r *registry) GetByName(name string) (*Entry, bool) {
	parts := strings.Split(name, ".")
	out := parts[len(parts)-1]
	lower := strings.ToLower(out)

	entry, ok := r.names[lower]
	return entry, ok
}

func NewRegistry() Registry {
	return &registry{
		names: make(map[string]*Entry),
	}
}
