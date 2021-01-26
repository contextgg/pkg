package es

import (
	"context"
	"time"

	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/metadata"
	"github.com/contextgg/pkg/types"
)

// NewEvent will create an event from data
func NewEvent(ctx context.Context, entity Entity, version int, data interface{}) events.Event {
	_, typeName := types.GetTypeName(data)
	timestamp := time.Now()
	meta := metadata.FromContext(ctx)

	return events.Event{
		Type:          typeName,
		Timestamp:     timestamp,
		AggregateID:   entity.GetID(),
		AggregateType: entity.GetTypeName(),
		Version:       version,
		Data:          data,
		Metadata:      meta,
	}
}
