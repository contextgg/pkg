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
	Upsert(obj interface{}, replaced ...string) *Entry
	GetByName(name string) (*Entry, bool)
	GetAll() []*Entry
}

type registry struct {
	names map[string]*Entry
}

func (r *registry) Upsert(obj interface{}, previousNames ...string) *Entry {
	t := GetElemType(obj)
	n := GetShortName(t)

	lower := strings.ToLower(n)
	e, ok := r.names[lower]
	if !ok {
		e = &Entry{
			Name:    n,
			Factory: TypeFactory(t),
		}
	}

	e.PreviousNames = append(e.PreviousNames, previousNames...)

	r.names[lower] = e
	for _, n := range previousNames {
		lower := strings.ToLower(n)
		r.names[lower] = e
	}
	return e
}

func (r *registry) GetByName(name string) (*Entry, bool) {
	parts := strings.Split(name, ".")
	out := parts[len(parts)-1]
	lower := strings.ToLower(out)

	entry, ok := r.names[lower]
	return entry, ok
}

func (r *registry) GetAll() []*Entry {
	var out []*Entry
	for _, e := range r.names {
		out = append(out, e)
	}
	return out
}

func NewRegistry() Registry {
	return &registry{
		names: make(map[string]*Entry),
	}
}
