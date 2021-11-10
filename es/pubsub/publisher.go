package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
)

type publisher struct {
	log    logger.Logger
	client *pubsub.Client
	topic  *pubsub.Topic
	codec  events.EventCodec
}

// NewPublisher creates a publisher and subscribes
func NewPublisher(log logger.Logger, cli *pubsub.Client, topicName string) (es.EventPublisher, error) {
	ctx := context.Background()

	topic, err := getTopic(ctx, log, cli, topicName)
	if err != nil {
		log.Error("getTopic", "err", err, "topicName", topicName)
		return nil, fmt.Errorf("getTopic: %v", err)
	}

	p := &publisher{
		log:    log,
		client: cli,
		topic:  topic,
		codec:  &json.EventCodec{},
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

	publishCtx := context.Background()
	res := c.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := res.Get(publishCtx); err != nil {
		c.log.Error("Could not publish event", "err", err)
		return err
	}

	c.log.Debug("Event Published via GCP pub/sub", "topic_id", c.topic.ID(), "event_type", event.Type, "event_aggregate_id", event.AggregateId, "event_aggregate_type", event.AggregateType)
	return nil
}
