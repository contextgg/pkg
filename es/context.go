package es

import (
	"context"
	"fmt"
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

func GetUnit(ctx context.Context) (Unit, error) {
	unit := UnitFromContext(ctx)
	if unit == nil {
		return nil, fmt.Errorf("no unit in context")
	}
	return unit, nil
}
