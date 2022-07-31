package es

import (
	"context"
)

type Key int

const (
	UnitKey Key = iota
)

func SetUnit(ctx context.Context, unit Unit) context.Context {
	return context.WithValue(ctx, UnitKey, unit)
}
func UnitFromContext(ctx context.Context) Unit {
	unit, ok := ctx.Value(UnitKey).(Unit)
	if ok {
		return unit
	}
	return nil
}
