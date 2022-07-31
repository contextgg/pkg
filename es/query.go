package es

import (
	"context"

	"github.com/contextgg/pkg/types"
)

type Query[T any] interface {
	Load(ctx context.Context, aggregateId string) (*T, error)
}

type query[T any] struct {
	entityName string
}

func (q *query[T]) Load(ctx context.Context, aggregateId string) (*T, error) {
	var item T
	if err := q.unit.Load(ctx, aggregateId, q.entityName, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func NewQuery[T any]() Query[T] {
	var item T
	entityName := types.GetTypeName(item)

	return &query[T]{
		entityName: entityName,
	}
}
