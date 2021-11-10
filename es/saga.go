package es

import (
	"context"

	"github.com/contextgg/pkg/events"
)

// SagaFunc for creating an saga
type SagaFunc func(Unit) Saga

// Saga takes a events and may return new commands
type Saga interface {
	Run(ctx context.Context, evt events.Event) ([]Command, error)
}
