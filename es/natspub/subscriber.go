package natspub

import (
	"context"
	"fmt"
	"time"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/events/json"
	"github.com/contextgg/pkg/logger"
	"github.com/nats-io/nats.go"
)

type subscriber struct {
	log          logger.Logger
	eventHandler es.EventHandler

	nc            *nats.Conn
	topicName     string
	subscriptions []*nats.Subscription

	codec events.EventCodec
	errCh chan es.EventBusError
}

func (c *subscriber) Errors() <-chan es.EventBusError {
	return c.errCh
}

func (c *subscriber) Start() {
	for _, sub := range c.subscriptions {
		go c.handle(sub)
	}
}

func (c *subscriber) handle(sub *nats.Subscription) {
	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err == nats.ErrTimeout {
			continue
		}

		if err != nil {
			c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not receive: %v", err)}
		}

		if msg != nil {
			c.handler(msg)
		}
	}
}

func (c *subscriber) handler(msg *nats.Msg) {
	ctx := context.Background()
	evt, ctx, err := c.codec.UnmarshalEvent(ctx, msg.Data)
	if err != nil {
		c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not unmarshal event: %s", err.Error()), Ctx: ctx}
		return
	}

	ctx = es.SetIsPublisher(ctx)

	// Notify all observers about the event.
	if err := c.eventHandler.HandleEvent(ctx, *evt); err != nil {
		c.errCh <- es.EventBusError{Err: fmt.Errorf("Could not handle event: %s", err.Error()), Ctx: ctx, Event: evt}
		return
	}
}

// NewSubscriber creates a subscriber
func NewSubscriber(log logger.Logger, eventHandler es.EventHandler, nc *nats.Conn, appId string, topicName string, listeners ...string) (es.EventSubscriber, error) {
	s := &subscriber{
		log:          log,
		eventHandler: eventHandler,
		nc:           nc,
		topicName:    topicName,
		codec:        &json.EventCodec{},
		errCh:        make(chan es.EventBusError, 1),
	}

	for _, l := range listeners {
		if len(l) == 0 {
			continue
		}
		if l == topicName {
			return nil, fmt.Errorf("Can not subscribe to the publishing topic: %s", topicName)
		}

		// create the subscription
		sub, err := nc.QueueSubscribeSync(l, appId)
		if err != nil {
			nc.Close()
			return nil, err
		}

		s.subscriptions = append(s.subscriptions, sub)
	}

	return s, nil
}
