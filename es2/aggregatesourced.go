package es2

import (
	"context"
	"fmt"
)

var ErrNotHandled = fmt.Errorf("Command not handled")

type AggregateSourced interface {
	StoreEventData(ctx context.Context, cmd Command) error
}

type BaseAggregateSourced struct {
}

func (a *BaseAggregateSourced) StoreEventData(ctx context.Context, cmd Command) error {
	// v := a.Version + len(a.events) + 1

	// e := NewEvent(ctx, a, v, data)
	// a.events = append(a.events, e)

	return nil
}

func NewBaseAggregateSourced(name string) *BaseAggregateSourced {
	return &BaseAggregateSourced{}
}
