package json

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/types"
)

type EventWithContext struct {
	AggregateId   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	Version       int                    `json:"version"`
	Type          string                 `json:"type"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          interface{}            `json:"data"`
	Metadata      map[string]interface{} `json:"metadata"`
	Context       map[string]string      `json:"context"`
}

// EventCodec is a codec for marshaling and unmarshaling events
// to and from bytes in JSON format.
type EventCodec struct{}

// MarshalEvent marshals an event into bytes in JSON format.
func (c *EventCodec) MarshalEvent(ctx context.Context, event *events.Event) ([]byte, error) {
	e := EventWithContext{
		AggregateId:   event.AggregateId,
		AggregateType: event.AggregateType,
		Version:       event.Version,
		Type:          event.Type,
		Timestamp:     event.Timestamp,
		Data:          event.Data,
		Metadata:      event.Metadata,
		Context:       events.MarshalContext(ctx),
	}

	// Marshal the event (using JSON for now).
	b, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not marshal event: %w", err)
	}

	return b, nil
}

// UnmarshalEvent unmarshals an event and supported parts of context from bytes.
func (c *EventCodec) UnmarshalEvent(ctx context.Context, b []byte) (*events.Event, context.Context, error) {
	out := struct {
		*EventWithContext

		Data json.RawMessage `json:"data"`
	}{}

	if err := json.Unmarshal(b, &out); err != nil {
		return nil, nil, fmt.Errorf("Could not decode event: %w", err)
	}

	typeData, ok := types.GetTypeData(out.Type)
	if !ok {
		return nil, nil, fmt.Errorf("Could not find type with name %s", out.Type)
	}

	data := typeData.Factory()
	if err := json.Unmarshal(out.Data, data); err != nil {
		return nil, nil, fmt.Errorf("Could not decode event data %w", err)
	}

	evt := &events.Event{
		AggregateId:   out.AggregateId,
		AggregateType: out.AggregateType,
		Version:       out.Version,
		Type:          out.Type,
		Timestamp:     out.Timestamp,
		Data:          data,
		Metadata:      out.Metadata,
	}

	if evt.Metadata == nil {
		evt.Metadata = make(map[string]interface{})
	}

	ctx = events.UnmarshalContext(ctx, out.Context)
	return evt, ctx, nil
}
