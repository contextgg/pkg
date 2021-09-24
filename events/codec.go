package events

import "context"

type EventCodec interface {
	// MarshalEvent marshals an event and the supported parts of context into bytes.
	MarshalEvent(context.Context, *Event) ([]byte, error)

	// UnmarshalEvent unmarshals an event and supported parts of context from bytes.
	UnmarshalEvent(context.Context, []byte) (*Event, context.Context, error)
}
