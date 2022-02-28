package pubsub

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
)

type subscriber struct {
	log           logger.Logger
	client        *pubsub.Client
	subscriptions []*pubsub.Subscription
	topic         *pubsub.Topic
	eventHandler  es.EventHandler
	errCh         chan es.EventBusError
	codec         events.EventCodec
}

func (c *subscriber) Errors() <-chan es.EventBusError {
	return c.errCh
}

func (c *subscriber) Start() {
	for _, sub := range c.subscriptions {
		go c.handle(sub)
	}
}

func (c *subscriber) handle(sub *pubsub.Subscription) {
	for {
		ctx := context.Background()
		if err := sub.Receive(ctx, c.handler); err != context.Canceled {
			select {
			case c.errCh <- es.EventBusError{Ctx: ctx, Err: fmt.Errorf("Could not receive: %v", err)}:
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

func (c *subscriber) handler(ctx context.Context, msg *pubsub.Message) {
	evt, ctx, err := c.codec.UnmarshalEvent(ctx, msg.Data)
	if err != nil && err == events.ErrTypeNotFound {
		msg.Ack()
		return
	}

	if err != nil {
		select {
		case c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not unmarshal event: %s", err.Error()), Ctx: ctx}:
		default:
		}
		msg.Nack()
		return
	}

	// Notify all observers about the event.
	if err := c.eventHandler.HandleEvent(ctx, *evt); err != nil {
		select {
		case c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not handle event: %s", err.Error()), Ctx: ctx, Event: evt}:
		default:
		}
		msg.Nack()
		return
	}

	msg.Ack()
}

func (c *subscriber) subscription(ctx context.Context, appId string, topicName string) error {
	sub, err := getSubscription(ctx, c.log, c.client, appId, topicName)
	if err != nil {
		return err
	}

	c.subscriptions = append(c.subscriptions, sub)
	return nil
}

// NewSubscriber creates a new subscriber
func NewSubscriber(log logger.Logger, eventHandler es.EventHandler, cli *pubsub.Client, appId string, topicName string, listeners ...string) (es.EventSubscriber, error) {
	ctx := context.Background()

	topic, err := getTopic(ctx, log, cli, topicName)
	if err != nil {
		log.Error("getTopic", "err", err, "topicName", topicName)
		return nil, fmt.Errorf("getTopic: %v", err)
	}

	p := &subscriber{
		log:          log,
		client:       cli,
		topic:        topic,
		eventHandler: eventHandler,
		errCh:        make(chan es.EventBusError, 1),
		codec:        &json.EventCodec{},
	}

	for _, l := range listeners {
		if len(l) == 0 {
			continue
		}
		if l == topicName {
			return nil, fmt.Errorf("Can not subscribe to the publishing topic: %s", topicName)
		}

		if err := p.subscription(ctx, appId, l); err != nil {
			return nil, err
		}
	}
	return p, nil
}
