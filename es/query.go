package es

import (
	"context"
	"fmt"

	"github.com/contextgg/pkg/types"
)

type Query[T Entity] interface {
	Load(ctx context.Context, id string) (T, error)
	Save(ctx context.Context, entities ...T) error
}

type query[T Entity] struct {
	name string
}

func (q *query[T]) Load(ctx context.Context, id string) (T, error) {
	var item T

	unit, err := GetUnit(ctx)
	if err != nil {
		return item, err
	}

	out, err := unit.Load(ctx, item, id)
	if err != nil {
		return item, err
	}

	result, ok := out.(T)
	if !ok {
		return item, fmt.Errorf("unexpected type: %T", out)
	}
	return result, nil
}

func (q *query[T]) Save(ctx context.Context, entities ...T) error {
	unit, err := GetUnit(ctx)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		if err := unit.Save(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

func NewQuery[T Entity]() Query[T] {
	var item T
	name := types.GetTypeName(item)

	return &query[T]{
		name: name,
	}
}
