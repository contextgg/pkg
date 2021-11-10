package natspub

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
)

type publisher struct {
	log logger.Logger

	nc        *nats.Conn
	topicName string
	codec     events.EventCodec
}

// NewPublisher creates a publisher and subscribes
func NewPublisher(log logger.Logger, nc *nats.Conn, topicName string) (es.EventPublisher, error) {
	p := &publisher{
		log:       log,
		nc:        nc,
		topicName: topicName,
		codec:     &json.EventCodec{},
	}

	return p, nil
}

// PublishEvent via pubsub
func (c *publisher) PublishEvent(ctx context.Context, event events.Event) error {
	data, err := c.codec.MarshalEvent(ctx, &event)
	if err != nil {
		c.log.Error("json.Marshal", "err", err)
		return err
	}

	if err := c.nc.Publish(c.topicName, data); err != nil {
		c.log.Error("Could not publish event", "err", err)
		return err
	}

	c.log.Debug("Event Published via GCP pub/sub", "topic_name", c.topicName, "event_type", event.Type, "event_aggregate_id", event.AggregateId, "event_aggregate_type", event.AggregateType)
	return nil
}
